package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"redgate.com/b/db"
	"redgate.com/b/db/sqlc"
	"redgate.com/b/handlers"
	"redgate.com/b/token"
)

func main() {
	l := log.New(os.Stdout, "RED-GATE-SERVER-", log.LstdFlags)
	ctx := context.Background()
	// load env
	err := godotenv.Load("dev.env")
	if err != nil {
		l.Fatalf("Error reding the .env %s", err)
	}

	// CRUD
	db, queries := db.Instantiate(l)
	if db == nil || queries == nil {
		l.Println("Exiting due to database connection error")
		return
	}
	defer db.Close()

	server := &http.Server{
		Addr:        "127.0.0.1:" + os.Getenv("PORT"),
		Handler:     defineMultiplexer(l, queries),
		IdleTimeout: 30 * time.Second,
		ReadTimeout: time.Second,
	}

	// now the startServer is run by a routine
	go startServer(server, l)

	// inorder to block the routine, we might use a channel (we can use wait group also)
	shut := make(chan os.Signal, 1)
	signal.Notify(shut, syscall.SIGINT, syscall.SIGTERM)

	<-shut // Block until a signal is received

	timeout_ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	stopServer(server, l, &timeout_ctx, &cancel)
}

func startServer(s *http.Server, l *log.Logger) {
	l.Println("🔥 Server is starting on", s.Addr)

	err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		l.Fatalln("Server is failed due to", err)
	}
}

func stopServer(s *http.Server, l *log.Logger, ctx *context.Context, cancel *context.CancelFunc) {
	l.Println("💅 Shutting down the server")
	s.Shutdown(*ctx)
	c := *cancel
	c()
}

func defineMultiplexer(l *log.Logger, q *sqlc.Queries) http.Handler {
	var u handlers.AuthedUser

	// reference to the handler
	hello_handler := handlers.NewHello(l)
	account_handler := handlers.NewAccountHandler(l, q, &u)
	token, err := token.NewPasetoMaker(os.Getenv("PASETO_KEY"))
	if err != nil {
		log.Fatal("Failed creating Paseto token")
	}
	auth_handler := handlers.NewAuthHandler(l, q, &u, &token)
	token_handler := handlers.NewTokenHandler(l, q, &u, &token)
	plate_handler := handlers.NewPlateIDHandler(l, q, &u)

	// handle multiplexer
	mux := http.NewServeMux()
	mux.Handle("/", hello_handler)
	mux.HandleFunc("/account/list", account_handler.ListAccountsH)

	mux.HandleFunc("/auth/login", auth_handler.Login)
	mux.HandleFunc("/auth/signup", auth_handler.Signup)
	mux.HandleFunc("/auth/renewToken", token_handler.RenewToken)

	mux.HandleFunc("/plate/create", plate_handler.CreatePlateHandler)
	mux.HandleFunc("/plate/verify", plate_handler.VerifyPlateHandler)

	corsMiddleware := cors.AllowAll().Handler

	handler := corsMiddleware(mux)

	return handler
}