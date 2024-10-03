package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"crproductos/db"
	"crproductos/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"strconv"
)

func main() {
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
