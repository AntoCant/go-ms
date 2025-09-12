package application

import (
	"github.com/google/uuid"

	"go-ms/internal/common/errs"
	"go-ms/internal/core"
	"go-ms/internal/core/domain"
)

type productService struct {
	repo core.ProductRepository
}

func NewProductService(r core.ProductRepository) core.ProductUseCase {
	return &productService{repo: r}
}

func (s *productService) CreateProduct(name string, price float64, stock int) (domain.Product, error) {
	if name == "" || price <= 0 {
		return domain.Product{}, errs.ErrBadRequest
	}
	p := domain.Product{
		Id:    uuid.NewString(),
		Name:  name,
		Price: price,
		Stock: stock,
	}
	if err := s.repo.Save(p); err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

func (s *productService) GetAllProducts() ([]domain.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) GetProductById(productId string) (domain.Product, error) {

	if productId == "" {
		return domain.Product{}, errs.PorductIdNotFound
	}
	return s.repo.FindByIdProduct(productId)
}

func (s *productService) UpdateProduct(productId string, name string, price float64, stock int) (domain.Product, error) {

	if productId == "" {
		return domain.Product{}, errs.PorductNotFound // o un ErrInvalidID/BadRequest si tienes
	}
	if price <= 0 {
		return domain.Product{}, errs.ErrBadRequest
	}
	if stock < 0 {
		return domain.Product{}, errs.ErrBadRequest
	}

	product, err := s.repo.FindByIdProduct(productId)
	if err != nil {
		return domain.Product{}, err
	}

	product.Name = name
	product.Price = price
	product.Stock = stock

	if err := s.repo.Save(product); err != nil {
		return domain.Product{}, err
	}

	return product, nil

}

func (s *productService) DeleteProduct(productId string) error {
	if productId == "" {
		return errs.PorductNotFound
	}

	if _, err := s.repo.FindByIdProduct(productId); err != nil {
		return err
	}
	return s.repo.DeleteProductById(productId)
}
