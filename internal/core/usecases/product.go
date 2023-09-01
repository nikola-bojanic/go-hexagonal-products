package usecases

import (
	"context"

	domain "github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/pkg/errors"
)

var _ ports.ProductUsecase = (*ProductService)(nil)

type ProductService struct {
	productRepo *repo.ProductRepository
}

func NewProductService(productRepo *repo.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) GetAllProducts(ctx context.Context) (*[]domain.Product, error) {
	products, err := s.productRepo.GetAllProducts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve products")
	}
	return products, nil
}
func (s *ProductService) FindProductById(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := s.productRepo.FindProductById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve a product")
	}
	return product, nil
}
func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (int64, error) {
	id, err := s.productRepo.InsertProduct(ctx, product)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to create a product")
	}
	return id, nil
}
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) (int64, error) {
	id, err := s.productRepo.DeleteProduct(ctx, id)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to delete a product")
	}
	return id, nil
}
func (s *ProductService) UpdateProduct(ctx context.Context, product *domain.Product, id int64) (int64, error) {
	id, err := s.productRepo.UpdateProduct(ctx, product, id)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to edit a product")
	}
	return id, nil
}
