package handlers

import (
	"log"
	"net/http"

	"redgate.com/b/db/sqlc"
	"redgate.com/b/utils"
)

func NewAccountHandler(l *log.Logger, q *sqlc.Queries, u *AuthedUser) *AccountHandler {
	var c uint = 0
	return &AccountHandler{&Handler{l, q, &c, u}}
}

func (ah *AccountHandler) CreateAccountH(w http.ResponseWriter, r *http.Request) {
	hp := HandlerParam{w, r, http.MethodPost, ah.createAccount}
	ah.h.handleRequest(hp, ah.h.u)
}

//	func (ah *AccountHandler) GetAccountH(w http.ResponseWriter, r *http.Request) {
//		hp := HandlerParam{w, r, http.MethodGet, ah.getAccount}
//		ah.h.handleRequest(hp, ah.h.u)
//	}
func (ah *AccountHandler) ListAccountsH(w http.ResponseWriter, r *http.Request) {
	hp := HandlerParam{w, r, http.MethodGet, ah.listAccounts}
	ah.h.handleRequest(hp, ah.h.u)
}

// func (ah *AccountHandler) DeleteAccountH(w http.ResponseWriter, r *http.Request) {
// 	hp := HandlerParam{w, r, http.MethodDelete, ah.deleteAccount}
// 	ah.h.handleRequest(hp, ah.h.u)
// }

// the implementation

func (ah *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) error {
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

	// Create accountParams using retrieved form values
	accountParams := sqlc.CreateAccountParams{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword, // Don't forget to hash the password
	}

	account, err := ah.h.q.CreateAccount(r.Context(), accountParams)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	toJSON(w, account)
	return nil
}

// func (ah *AccountHandler) getAccount(w http.ResponseWriter, r *http.Request) error {
// 	accountID := r.URL.Query().Get("account_id")

// 	// if accountID != int(ah.h.u.UserID) {
// 	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
// 	// 	return errors.New("unauthorized")
// 	// }

// 	account, err := ah.h.q.GetAccount(r.Context(), accountID)
// 	if err != nil {
// 		http.Error(w, "Account not found", http.StatusNotFound)
// 		return err
// 	}

// 	toJSON(w, account)
// 	return nil
// }

func (ah *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := ah.h.q.ListAccounts(r.Context())
	if err != nil {
		http.Error(w, "Error listing accounts", http.StatusInternalServerError)
		return err
	}

	// TO DO : only admin can use this query
	// if r.Header.Get("admin") == "false" {
	// 	http.Error(w, "unauthorized", http.StatusNonAuthoritativeInfo)
	// 	return errors.New("unauthorized")
	// }

	toJSON(w, accounts)
	return nil
}
