package service

import (
	"crproductos/internal/models"
	"crproductos/internal/repository"
)

type productService struct {
	repo repository.ProductRepository
}
type ProductService interface {
	GetAllProducts() ([]models.ProductResponse, error)
	GetProductById(id string) (models.Product, error)
	CreateProduct(product models.ProductResponse) (models.ProductResponse, error)
	DeleteProduct(id string) error
	UpdateProduct(id string, product models.ProductResponse) (models.ProductResponse, error)
	PatchProduct(id string, product models.ProductResponse) (models.Product, error)
	PatchStore(id string, jsonStore []byte) (models.Product, error)
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}
func (s *productService) GetAllProducts() ([]models.ProductResponse, error) {
	return s.repo.GetAllProducts()
}

func (s *productService) GetProductById(id string) (models.Product, error) {
	return s.repo.GetProductById(id)
}

func (s *productService) CreateProduct(product models.ProductResponse) (models.ProductResponse, error) {
	return s.repo.CreateProduct(product)
}
func (s *productService) DeleteProduct(id string) error {
	return s.repo.DeleteProduct(id)
}
func (s *productService) UpdateProduct(id string, product models.ProductResponse) (models.ProductResponse, error) {
	return s.repo.UpdateProduct(id, product)
}
func (s *productService) PatchProduct(id string, product models.ProductResponse) (models.Product, error) {
	return s.repo.PatchProduct(id, product)
}
func (s *productService) PatchStore(id string, jsonStore []byte) (models.Product, error) {
	return s.repo.PatchStore(id, jsonStore)
}
