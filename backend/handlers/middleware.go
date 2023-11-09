package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"redgate.com/b/token"
)

func (h *Handler) handleRequest(hp HandlerParam, u *AuthedUser) {
	*h.c++
	var err error = nil
	defer func() {
		apiLog(h.l, h.c, &hp.r.RequestURI, err)
	}()

	err = checkHTTPMethod(hp.w, hp.r.Method, hp.method)
	if err != nil {
		return
	}

	err = checkAuthorization(hp.w, hp.r, u)
	if err != nil {
		return
	}

	err = hp.handlerFunc(hp.w, hp.r)
	if err != nil {
		return
	}
}

func apiLog(l *log.Logger, counter *uint, url *string, err error) {
	var status string
	if err == nil {
		status = "SUCCESS"
	} else {
		status = err.Error()
	}

	l.Printf("[%d] [%s] [%s]", *counter, *url, status)
}

func toJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func checkHTTPMethod(w http.ResponseWriter, reqMethod, desMethod string) error {
	if reqMethod != desMethod {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return errors.New("invalid http method")
	}
	return nil
}

func checkAuthorization(w http.ResponseWriter, r *http.Request, u *AuthedUser) error {
	if strings.HasPrefix(r.URL.Path, "/auth/") {
		// If the request URI starts with /auth/, skip authorization
		return nil
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "No authorization header found in the request", http.StatusNonAuthoritativeInfo)
		return errors.New("unauthorized")
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Only accept bearer token", http.StatusNonAuthoritativeInfo)
		return errors.New("invalid token format")
	}

	bearer_token := authHeader[len("Bearer "):]
	PASETO_KEY := os.Getenv("PASETO_KEY")
	maker, err := token.NewPasetoMaker(PASETO_KEY)
	if err != nil {
		http.Error(w, "Failed to verify", http.StatusInternalServerError)
		return errors.New("cannot create paseto")
	}

	payload, err := maker.VerifyToken(bearer_token)
	if err != nil {
		http.Error(w, "Invalid Token or Expired", http.StatusNonAuthoritativeInfo)
		return errors.New("invalid token or expired")
	}

	u.UserID = payload.AccountID
	u.Username = payload.Username

	return nil
}