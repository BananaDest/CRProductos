package handlers

import (
	"crproductos/db"
	"crproductos/models"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type Server struct {
	Router *chi.Mux
	DB     *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{Router: chi.NewRouter(), DB: db}
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(render.SetContentType(render.ContentTypeJSON))
	s.Router.Get("/", rootHandler)
	s.Router.Route("/products", func(r chi.Router) {
		r.Get("/", s.GetAllProducts)
		r.Get("/{id}", s.GetProductById)
		r.Post("/", s.CreateProduct)
		r.Put("/{id}", s.UpdateProduct)
		r.Delete("/{id}", s.DeleteProduct)
		r.Patch("/{id}", s.PatchProduct)
		r.Patch("/{id}/store", s.PatchStore)
	})
}

func (s *Server) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := db.GetAllProducts(s.DB)
	if err != nil {
		http.Error(w, "Failed getting all products", http.StatusInternalServerError)
	}
	render.JSON(w, r, products)

}
func (s *Server) GetProductById(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	product, err := db.GetProductById(s.DB, id)
	if err != nil {
		http.Error(w, "Failed getting product", http.StatusInternalServerError)
	}
	render.JSON(w, r, product.ToJSON())
}
func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	product, err := db.CreateProduct(s.DB, product)
	if err != nil {
		http.Error(w, "Failed creating product", http.StatusInternalServerError)
	}
	render.JSON(w, r, product)
}
func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	if err := db.DeleteProduct(s.DB, id); err != nil {
		http.Error(w, "Failed deleting product", http.StatusInternalServerError)
	}
	w.Write([]byte("Delete successful"))
}

func (s *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var product models.ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	product, err := db.UpdateProduct(s.DB, id, product)
	if err != nil {
		http.Error(w, "Failed updating product", http.StatusInternalServerError)
	}

	render.JSON(w, r, product)

}
func (s *Server) PatchProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var product models.ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	updatedProduct, err := db.PatchProduct(s.DB, id, product)
	if err != nil {
		http.Error(w, "Failed patching product", http.StatusInternalServerError)
	}

	render.JSON(w, r, updatedProduct.ToJSON())
}

func (s *Server) PatchStore(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var store models.Stores
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	jsonStore, err := json.Marshal(store)
	if err != nil {
		http.Error(w, "Unable to convert data to JSON", http.StatusInternalServerError)
		return
	}
	updatedProduct, err := db.PatchStore(s.DB, id, jsonStore)
	if err != nil {
		http.Error(w, "Failed patching store", http.StatusInternalServerError)
	}
	render.JSON(w, r, updatedProduct.ToJSON())
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
