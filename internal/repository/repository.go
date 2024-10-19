package repository

import "crproductos/internal/models"

type ProductRepository interface {
	GetAllProducts() ([]models.ProductResponse, error)
	GetProductById(id string) (models.Product, error)
	CreateProduct(product models.ProductResponse) (models.ProductResponse, error)
	DeleteProduct(id string) error
	UpdateProduct(id string, product models.ProductResponse) (models.ProductResponse, error)
	PatchProduct(id string, product models.ProductResponse) (models.Product, error)
	PatchStore(id string, jsonStore []byte) (models.Product, error)
}
