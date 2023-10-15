package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// TODO adding sqlc for the ORM
func Instantiate(l *log.Logger) (*sql.DB, interface{}) {
	err := godotenv.Load("dev.env")
	if err != nil {
		l.Fatalf("Error reding the .env %s", err)
	}

	db, err1 := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_CONNECTION_STRING"))
	if err1 != nil {
		l.Println("Error creating DB connection", err1)
		return nil, "nil"
	}

	err2 := db.Ping()
	if err2 != nil {
		l.Println("Error connecting to DB ", err2)
		return nil, "nil"
	}

	l.Println("🛢️  DB Connected")
	return db, "nil"
}
