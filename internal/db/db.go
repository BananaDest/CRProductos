// Package db handles all logic related to interacting with the database
// The package exports CRUD functions for the Product table in postgres
// Used mainly within http handlers to allow them to interact with the DB
package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

// ConnectToPostgres: Loads environment variables and then creates the connection string for postgres
// Returns the sql DB object from database/Sql
func ConnectToPostgres() *sql.DB {
	err := godotenv.Load("config.development.env")
	if err != nil {
		log.Fatalf("Error loading .env files: %v", err)
	}
	dbhost := os.Getenv("DB_HOST")
	dbportStr := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbport, err := strconv.Atoi(dbportStr)
	if err != nil {
		log.Fatalf("error converting: %v", err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", psqlInfo)
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
