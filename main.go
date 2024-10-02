package main

import (
	"fmt"
	"net/http"

	"crproductos/db"
	"crproductos/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
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
	db := db.ConnectToPostgres(psqlInfo)
	defer db.Close()
	server := handlers.NewServer(db)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", rootHandler)
	r.Route("/products", func(r chi.Router) {
		r.Get("/", server.GetAllProducts)
		r.Get("/{id}", server.GetProductById)
		r.Post("/", server.CreateProduct)
		r.Put("/{id}", server.UpdateProduct)
		r.Delete("/{id}", server.DeleteProduct)
		r.Patch("/{id}", server.PatchProduct)
	})
	http.ListenAndServe(":8080", r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
