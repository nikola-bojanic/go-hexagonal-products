package order

import (
	"time"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
)

type OrderModel struct {
	ID           string                 `json:"id"`
	ProductItems *[]OrderedProductModel `json:"product_items"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

type OrderedProductModel struct {
	ProductId int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

func (e *OrderModel) FromDomain(order *domain.Order) {
	if e == nil || order == nil {
		return
	}
	e.ID = order.ID
	e.Status = order.Status
	e.CreatedAt = order.CreatedAt
	e.UpdatedAt = order.UpdatedAt
	var products []OrderedProductModel = []OrderedProductModel{}
	for _, item := range *order.ProductItems {
		orderedProduct := OrderedProductModel{}
		orderedProduct.FromDomain(&item)
		products = append(products, orderedProduct)
	}
	e.ProductItems = &products
}

func (e *OrderModel) ToDomain() *domain.Order {
	if e == nil {
		return &domain.Order{}
	}
	var products []domain.OrderedProduct = []domain.OrderedProduct{}
	for _, item := range *e.ProductItems {
		product := item.ToDomain()
		products = append(products, *product)
	}
	return &domain.Order{
		ID:           e.ID,
		Status:       e.Status,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		ProductItems: &products,
	}
}

func (e *OrderedProductModel) FromDomain(orderedProduct *domain.OrderedProduct) {
	if e == nil || orderedProduct == nil {
		return
	}
	e.ProductId = orderedProduct.ProductId
	e.Quantity = orderedProduct.Quantity
}

func (e *OrderedProductModel) ToDomain() *domain.OrderedProduct {
	if e == nil {
		return &domain.OrderedProduct{}
	}
	return &domain.OrderedProduct{
		Quantity:  e.Quantity,
		ProductId: e.ProductId,
	}
}
