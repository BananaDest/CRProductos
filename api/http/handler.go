package http

import (
	"crproductos/internal/models"
	"crproductos/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{service: svc}
}
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		http.Error(w, "Failed getting all products", http.StatusInternalServerError)
	}
	render.JSON(w, r, products)

}
func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	product, err := h.service.GetProductById(id)
	if err != nil {
		http.Error(w, "Failed getting product", http.StatusInternalServerError)
	}
	render.JSON(w, r, product.ToJSON())
}
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	product, err := h.service.CreateProduct(product)
	if err != nil {
		http.Error(w, "Failed creating product", http.StatusInternalServerError)
	}
	render.JSON(w, r, product)
}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	if err := h.service.DeleteProduct(id); err != nil {
		http.Error(w, "Failed deleting product", http.StatusInternalServerError)
	}
	w.Write([]byte("Delete successful"))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var product models.ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	product, err := h.service.UpdateProduct(id, product)
	if err != nil {
		http.Error(w, "Failed updating product", http.StatusInternalServerError)
	}

	render.JSON(w, r, product)

}
func (h *ProductHandler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var product models.ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	updatedProduct, err := h.service.PatchProduct(id, product)
	if err != nil {
		http.Error(w, "Failed patching product", http.StatusInternalServerError)
	}

	render.JSON(w, r, updatedProduct.ToJSON())
}

func (h *ProductHandler) PatchStore(w http.ResponseWriter, r *http.Request) {
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
	updatedProduct, err := h.service.PatchStore(id, jsonStore)
	if err != nil {
		http.Error(w, "Failed patching store", http.StatusInternalServerError)
	}
	render.JSON(w, r, updatedProduct.ToJSON())
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
