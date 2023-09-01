package category

import (
	"time"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
)

type CategoryModel struct {
	Id        int       `json:"categoryId"`
	Name      string    `json:"categoryName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (e *CategoryModel) FromDomain(category *domain.Category) {
	if e == nil || category == nil {
		return
	}
	e.Id = category.Id
	e.Name = category.Name
	e.UpdatedAt = category.UpdatedAt
	e.CreatedAt = category.CreatedAt
}

func (e *CategoryModel) ToDomain() *domain.Category {
	if e == nil {
		return &domain.Category{}
	}

	return &domain.Category{
		Id:        e.Id,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
