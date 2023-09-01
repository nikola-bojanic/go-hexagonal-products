package product

import (
	"time"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/handlers/category"
)

type ProductModel struct {
	ID               int                     `json:"productId"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"shortDescription"`
	Description      string                  `json:"description"`
	Price            float32                 `json:"price"`
	Quantity         int                     `json:"quantity"`
	Category         *category.CategoryModel `json:"category"`
	CreatedAt        time.Time               `json:"createdAt"`
	UpdatedAt        time.Time               `json:"updatedAt"`
}

func (e *ProductModel) FromDomain(product *domain.Product) {
	if e == nil || product == nil {
		return
	}

	e.ID = product.ProductId
	e.Name = product.Name
	e.ShortDescription = product.ShortDescription
	e.Description = product.Description
	e.Price = product.Price
	e.Quantity = product.Quantity
	e.Category = &category.CategoryModel{}
	e.Category.FromDomain(product.Category)
	e.CreatedAt = product.CreatedAt
	e.UpdatedAt = product.UpdatedAt

}

func (e *ProductModel) ToDomain() *domain.Product {
	if e == nil {
		return &domain.Product{}
	}
	return &domain.Product{
		ProductId:        e.ID,
		Name:             e.Name,
		ShortDescription: e.ShortDescription,
		Description:      e.Description,
		Price:            e.Price,
		Quantity:         e.Quantity,
		Category:         e.Category.ToDomain(),
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}
