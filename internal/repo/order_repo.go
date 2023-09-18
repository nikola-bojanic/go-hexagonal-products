package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/database"
)

var _ ports.OrderRepo = (*OrderRepository)(nil)

type OrderRepository struct {
	db                     *database.DB
	OrderProductRepository *OrderProductRepository
	UserRepository         *UserRepository
}

func NewOrderRepository(db *database.DB) *OrderRepository {
	return &OrderRepository{
		db:                     db,
		OrderProductRepository: NewOrderProductRepository(db),
		UserRepository:         NewUserRepository(db),
	}
}

func (repo *OrderRepository) FindOrderById(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	var userId string
	err := repo.db.QueryRow(ctx, `SELECT id, status, user_id, created_at, updated_at FROM hex_fwk.order WHERE id = $1`, id).Scan(&order.ID, &order.Status, &userId, &order.CreatedAt, &order.UpdatedAt)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("order not found")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	productItems, err := repo.OrderProductRepository.GetProducts(ctx, id)
	if err != nil {
		return nil, err
	}
	order.ProductItems = productItems
	user, err := repo.UserRepository.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	order.User = user
	return &order, nil
}
func (repo *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	order.Status = "CREATED"
	err := repo.db.QueryRow(ctx, `INSERT INTO hex_fwk.order (status, user_id) VALUES ($1, $2) RETURNING id, status, user_id, created_at, updated_at`, order.Status, order.User.ID).
		Scan(&order.ID, &order.Status, &order.User.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	for _, product := range *order.ProductItems {
		err := repo.OrderProductRepository.Add(ctx, order.ID, product.ProductId, product.Quantity)
		if err != nil {
			return nil, err
		}
	}
	productItems, err := repo.OrderProductRepository.GetProducts(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	user, err := repo.UserRepository.FindByID(ctx, order.User.ID)
	if err != nil {
		return nil, err
	}
	order.User = user
	order.ProductItems = productItems
	return order, nil
}
func (repo *OrderRepository) UpdateOrderStatus(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	updatedAt := time.Now()
	_, err := repo.db.Exec(ctx, `UPDATE hex_fwk.order SET status = $2, updated_at = $3 WHERE id = $1`, order.ID, order.Status, updatedAt)
	if err != nil {
		return nil, err
	}
	order.UpdatedAt = updatedAt
	return order, nil
}

func (repo *OrderRepository) DeleteOrder(ctx context.Context, order *domain.Order) error {
	for _, product := range *order.ProductItems {
		err := repo.OrderProductRepository.Delete(ctx, order.ID, product.ProductId)
		if err != nil {
			return err
		}
	}
	_, err := repo.db.Exec(ctx, `DELETE FROM hex_fwk.order WHERE id = $1`, order.ID)
	if err != nil {
		return err
	}
	return nil
}
