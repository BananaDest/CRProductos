package main

import (
	apiHttp "crproductos/api/http"
	"crproductos/internal/db"
	"crproductos/internal/repository"
	"crproductos/internal/service"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// Self explanatory, need to look if there is way to mock db to separate tests into unit and integration testing
	db := db.ConnectToPostgres()
	defer db.Close()
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := apiHttp.NewProductHandler(productService)
	server := apiHttp.NewServer()
	server.MountHandlers(productHandler)
	http.ListenAndServe(":8080", server.Router)
}
