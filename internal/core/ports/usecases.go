package ports

import (
	"context"

	domain "github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type ProductUsecase interface {
	GetAllProducts(ctx context.Context) (*[]domain.Product, error)
	FindProductById(ctx context.Context, id int64) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) (int64, error)
	DeleteProduct(ctx context.Context, id int64) (int64, error)
	UpdateProduct(ctx context.Context, product *domain.Product, id int64) (int64, error)
}

type CategoryUsecase interface {
	GetAllCategories(ctx context.Context) (*[]domain.Category, error)
	FindCategoryById(ctx context.Context, id int64) (*domain.Category, error)
	CreateCategory(ctx context.Context, category *domain.Category) (int64, error)
	DeleteCategory(ctx context.Context, id int64) (int64, error)
	UpdateCategory(ctx context.Context, category *domain.Category, id int64) (int64, error)
}

type OrderUsecase interface {
	FindOrderById(ctx context.Context, id string) (*domain.Order, error)
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, order *domain.Order) (*domain.Order, error)
	DeleteOrder(ctx context.Context, order *domain.Order) error
}
