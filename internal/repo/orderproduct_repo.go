package repo

import (
	"context"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/database"
)

var _ ports.OrderProductRepo = (*OrderProductRepository)(nil)

type OrderProductRepository struct {
	db                *database.DB
	ProductRepository *ProductRepository
}

func NewOrderProductRepository(db *database.DB) *OrderProductRepository {
	return &OrderProductRepository{
		db:                db,
		ProductRepository: NewProductRepository(db),
	}
}

func (repo *OrderProductRepository) GetProducts(ctx context.Context, orderId string) (*[]domain.OrderedProduct, error) {
	var products []domain.OrderedProduct
	rows, err := repo.db.Query(ctx, `SELECT product_id, quantity FROM hex_fwk.order_product WHERE order_id = $1`, orderId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var productId int64
		var quantity int
		err = rows.Scan(&productId, &quantity)
		if err != nil {
			return nil, err
		}
		product, err := repo.ProductRepository.FindProductById(ctx, productId)
		if err != nil {
			return nil, err
		}
		var orderProduct *domain.OrderedProduct
		orderProduct.ProductId = int64(product.ProductId)
		orderProduct.Quantity = quantity
		products = append(products, *orderProduct)
	}
	return &products, err
}

func (repo *OrderProductRepository) Add(ctx context.Context, orderId string, productId int64) error {
	_, err := repo.db.Exec(ctx, `INSERT INTO hex_fwk.order_product (order_id, product_id) VALUES ($1, $2)`, orderId, productId)
	if err != nil {
		return err
	}
	return nil
}
