package ports

import (
	"context"

	domain "github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
)

type UserRepo interface {
	Insert(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type ProductRepo interface {
	GetAllProducts(ctx context.Context) (*[]domain.Product, error)
	FindProductById(ctx context.Context, id int64) (*domain.Product, error)
	InsertProduct(ctx context.Context, product *domain.Product) (int64, error)
	DeleteProduct(ctx context.Context, id int64) (int64, error)
	UpdateProduct(ctx context.Context, product *domain.Product, id int64) (int64, error)
}

type CategoryRepo interface {
	GetAllCategories(ctx context.Context) (*[]domain.Category, error)
	FindCategoryById(ctx context.Context, id int64) (*domain.Category, error)
	InsertCategory(ctx context.Context, category *domain.Category) (int64, error)
	DeleteCategory(ctx context.Context, id int64) (int64, error)
	UpdateCategory(ctx context.Context, category *domain.Category, id int64) (int64, error)
}

type OrderRepo interface {
	FindOrderById(ctx context.Context, id string) (*domain.Order, error)
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, order *domain.Order) (*domain.Order, error)
}

type OrderProductRepo interface {
	Add(ctx context.Context, orderId string, productId int64, quantity int) error
	GetProducts(ctx context.Context, orderId string) (*[]domain.OrderedProduct, error)
}
