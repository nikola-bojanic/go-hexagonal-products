package usecases

import (
	"context"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/pkg/errors"
)

var _ ports.OrderUsecase = (*OrderService)(nil)

type OrderService struct {
	orderRepo   *repo.OrderRepository
	productRepo *repo.ProductRepository
}

func NewOrderService(orderRepo *repo.OrderRepository, productRepo *repo.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) FindOrderById(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.orderRepo.FindOrderById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve an product")
	}
	return order, nil
}
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	validStatus := order.Status != "" || order.Status != "CREATED" || order.Status != "PENDING" || order.Status != "COMPLETED" || order.Status != "CLOSED"
	if !validStatus {
		return nil, errors.New("invalid order status")
	}
	for _, item := range *order.ProductItems {
		product, err := s.productRepo.FindProductById(ctx, int64(item.ProductId))
		if err != nil {
			return nil, errors.New("product doesn't exist")
		}
		if product.Quantity < item.Quantity {
			return nil, errors.New("not enough items")
		}
	}
	created, err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create an order")
	}
	return created, nil
}
func (s *OrderService) UpdateOrderStatus(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	validStatus := order.Status == "" || order.Status != "CREATED" && order.Status != "PENDING" && order.Status != "COMPLETED" && order.Status != "CLOSED"
	if !validStatus {
		return nil, errors.New("invalid order status")
	}
	updated, err := s.orderRepo.UpdateOrderStatus(ctx, order)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update an order")
	}
	return updated, nil
}
