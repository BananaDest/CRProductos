package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "crproducts"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db := connectToPostgres(psqlInfo)
	defer db.Close()
	rows, err := db.Query("select * from crproducts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		//TODO: generate structure and add rows
		//TODO: test driven development
		fmt.Print(rows.Scan)
	}
	r := chi.NewRouter()
	r.Get("/", rootHandler)
	http.ListenAndServe(":8080", r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func connectToPostgres(connectionString string) *sql.DB {
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
