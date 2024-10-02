package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectToPostgres(connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect ", err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping postgres", err)
	}
	fmt.Println("connected to postgres")
	return db
}
