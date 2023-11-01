package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"redgate.com/b/db/sqlc"
	"redgate.com/b/token"
)

func NewTokenHandler(l *log.Logger, q *sqlc.Queries, u *AuthedUser, t *token.Maker) *AuthHandler {
	var c uint = 0
	return &AuthHandler{&Handler{l, q, &c, u}, *t}
}

func (auth_h *AuthHandler) RenewToken(w http.ResponseWriter, r *http.Request) {
	hp := HandlerParam{w, r, http.MethodPost, auth_h.renewToken}
	auth_h.h.handleRequest(hp, nil)
}

func (auth_h *AuthHandler) renewToken(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return err
	}

	// Retrieve form values
	refresh_token := r.FormValue("refresh_token")

	refresh_payload, err := auth_h.t.VerifyToken(refresh_token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return errors.New("invalid refresh token")
	}

	session, err := auth_h.h.q.GetSession(r.Context(), refresh_payload.ID)
	if err != nil {
		http.Error(w, "Cannot get Session", http.StatusInternalServerError)
		return errors.New("Cannot get Session")
	}

	if session.IsBlocked {
		http.Error(w, "Session is blocked", http.StatusUnauthorized)
		return errors.New("session is blocked")
	}

	if session.AccountID != refresh_payload.AccountID ||
		session.Username != refresh_payload.Username {
		http.Error(w, "identity not match", http.StatusUnauthorized)
		return errors.New("identity not match")
	}

	if session.RefreshToken != refresh_token {
		http.Error(w, "Mismatch token", http.StatusUnauthorized)
		return errors.New("mismatch token")
	}

	if time.Now().After(session.ExpiresAt) {
		http.Error(w, "Refresh token expire", http.StatusUnauthorized)
		return errors.New("refresh token expire")
	}

	access_token_duration := time.Minute * ACCESS_TOKEN_DURATION
	a_token, a_payload, err := auth_h.t.GenerateToken(session.AccountID, session.Username, access_token_duration)
	if err != nil {
		http.Error(w, "Failed to generate new token", http.StatusInternalServerError)
		return errors.New("failed generate new access token for user")
	}

	res := renewAccessTokenResponse{
		AccessToken:          a_token,
		AccessTokenExpiresAt: a_payload.ExpiredAt,
	}

	w.WriteHeader(http.StatusCreated)
	toJSON(w, res)
	return nil
}
