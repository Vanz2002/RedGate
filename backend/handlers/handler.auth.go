package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"redgate.com/b/db/sqlc"
	"redgate.com/b/token"
	"redgate.com/b/utils"
)

const (
	ACCESS_TOKEN_DURATION  = 100
	REFRESH_TOKEN_DURATION = 12
)

func NewAuthHandler(l *log.Logger, q *sqlc.Queries, u *AuthedUser, t *token.Maker) *AuthHandler {
	var c uint = 0
	return &AuthHandler{&Handler{l, q, &c, u}, *t}
}

func (auth_h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	hp := HandlerParam{w, r, http.MethodPost, auth_h.login}
	auth_h.h.handleRequest(hp, nil)
}

func (auth_h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	hp := HandlerParam{w, r, http.MethodPost, auth_h.signup}
	auth_h.h.handleRequest(hp, nil)
}

func (auth_h *AuthHandler) login(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return err
	}

	// Retrieve form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	ok := utils.EmailIsValid(email)
	if !ok {
		http.Error(w, "Invalid Email", http.StatusInternalServerError)
		return errors.New("invalid email")
	}

	ok = utils.PasswordIsValid(password)
	if !ok {
		http.Error(w, "Invalid Password", http.StatusInternalServerError)
		return errors.New("invalid password")
	}

	// check if user exist
	user, err := auth_h.h.q.GetAccountbyEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "Account not found! Register first", http.StatusInternalServerError)
		return errors.New("account not found, register first")
	}

	access_token_duration := time.Minute * ACCESS_TOKEN_DURATION
	a_token, a_payload, err := auth_h.t.GenerateToken(user.AccountID, user.Username, access_token_duration)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return errors.New("failed generate access token for user")
	}

	refresh_token_duration := time.Hour * REFRESH_TOKEN_DURATION
	r_token, r_payload, err := auth_h.t.GenerateToken(user.AccountID, user.Username, refresh_token_duration)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return errors.New("failed generate refresh token for user")
	}

	sessionParams := sqlc.CreateSessionParams{
		ID:           r_payload.ID,
		AccountID:    r_payload.AccountID,
		Username:     r_payload.Username,
		RefreshToken: r_token,
		ExpiresAt:    r_payload.ExpiredAt,
	}

	session, err := auth_h.h.q.CreateSession(r.Context(), sessionParams)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return errors.New("failed to create a session")
	}

	res := LoginUserResponse{
		SessionID:      session.ID,
		AccessToken:    a_token,
		AccessTokenEx:  a_payload.ExpiredAt,
		RefreshToken:   r_token,
		RefreshTokenEx: r_payload.ExpiredAt,
		UserID:         user.AccountID,
		Username:       user.Username,
	}

	w.WriteHeader(http.StatusCreated)
	toJSON(w, res)
	return nil
}

func (auth_h *AuthHandler) signup(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return err
	}

	// Retrieve form values
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashedPassword, _ := utils.HashPassword(password)

	ok := utils.EmailIsValid(email)
	if !ok {
		http.Error(w, "Invalid Email", http.StatusInternalServerError)
		return errors.New("invalid email")
	}

	ok = utils.PasswordIsValid(password)
	if !ok {
		http.Error(w, "Invalid Password", http.StatusInternalServerError)
		return errors.New("invalid password")
	}

	// Create accountParams using retrieved form values
	accountParams := sqlc.CreateAccountParams{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword, // Don't forget to hash the password
	}

	account, err := auth_h.h.q.CreateAccount(r.Context(), accountParams)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	toJSON(w, account)
	return nil
}
