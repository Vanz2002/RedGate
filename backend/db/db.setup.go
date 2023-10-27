package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"redgate.com/b/db/sqlc"
)

// TODO adding sqlc for the ORM
func Instantiate(l *log.Logger) (*sql.DB, *sqlc.Queries) {
	err := godotenv.Load("dev.env")
	if err != nil {
		l.Fatalf("Error reding the .env %s", err)
	}

	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")
	connection_string := "postgresql://" + DB_USER + ":" + DB_PASS + "@localhost:5432/" + DB_NAME + "?sslmode=disable"

	db, err1 := sql.Open(os.Getenv("DB_DRIVER"), connection_string)
	if err1 != nil {
		l.Println("Error creating DB connection", err1)
		return nil, nil
	}

	err2 := db.Ping()
	if err2 != nil {
		l.Println("Error connecting to DB ", err2)
		return nil, nil
	}

	l.Println("🛢️  DB Connected")
	return db, sqlc.New(db)
}
