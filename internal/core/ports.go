package core

import "go-ms/internal/core/domain"

// Inbound port (caso de uso)
type ProductUseCase interface {
	CreateProduct(name string, price float64, stock int) (domain.Product, error)
	GetAllProducts() ([]domain.Product, error)
	GetProductById(productId string) (domain.Product, error)
	UpdateProduct(productId string, name string, price float64, stock int) (domain.Product, error)
	DeleteProduct(productId string) error
}

// Outbound port (dependencia)
type ProductRepository interface {
	Save(domain.Product) error
	FindAll() ([]domain.Product, error)
	FindByIdProduct(productId string) (domain.Product, error)
	DeleteProductById(productId string) error
}
