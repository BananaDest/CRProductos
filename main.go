package main

import (
	"crproductos/db"
	"crproductos/handlers"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {

	db := db.ConnectToPostgres()
	defer db.Close()
	server := handlers.NewServer(db)
	server.MountHandlers()
	http.ListenAndServe(":8080", server.Router)
}
