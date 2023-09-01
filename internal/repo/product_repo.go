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

var _ ports.ProductRepo = (*ProductRepository)(nil)

type ProductRepository struct {
	db                 *database.DB
	CategoryRepository *CategoryRepository
}

func NewProductRepository(db *database.DB) *ProductRepository {
	return &ProductRepository{
		db:                 db,
		CategoryRepository: NewCategoryRepository(db),
	}
}
func (repo *ProductRepository) GetAllProducts(ctx context.Context) (*[]domain.Product, error) {
	var products []domain.Product
	rows, err := repo.db.Query(ctx, `SELECT id, name, short_description, description, price, category_id, quantity, created_at, updated_at FROM hex_fwk.product`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product domain.Product
		var categoryId int64
		err = rows.Scan(&product.ProductId, &product.Name, &product.ShortDescription, &product.Description, &product.Price,
			&categoryId, &product.Quantity, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		product.Category, err = repo.CategoryRepository.FindCategoryById(ctx, categoryId)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return &products, err
}

func (repo *ProductRepository) FindProductById(ctx context.Context, id int64) (*domain.Product, error) {
	var product domain.Product
	var categoryId int64
	err := repo.db.QueryRow(ctx, `SELECT id, name, short_description, description, price, category_id, quantity, created_at, updated_at FROM hex_fwk.product WHERE id = $1`, id).Scan(&product.ProductId, &product.Name, &product.ShortDescription, &product.Description, &product.Price,
		&categoryId, &product.Quantity, &product.CreatedAt, &product.UpdatedAt)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("product not found")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	product.Category, err = repo.CategoryRepository.FindCategoryById(ctx, categoryId)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *ProductRepository) InsertProduct(ctx context.Context, product *domain.Product) (int64, error) {
	var id int64
	err := repo.db.QueryRow(ctx, `INSERT INTO hex_fwk.product (name, short_description, description, price, quantity, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		product.Name, product.ShortDescription, product.Description, product.Price, product.Quantity, int64(product.Category.Id)).
		Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *ProductRepository) DeleteProduct(ctx context.Context, id int64) (int64, error) {
	res, err := repo.db.Exec(ctx, `DELETE FROM hex_fwk.product WHERE id = $1`,
		id)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (repo *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product, id int64) (int64, error) {
	updatedAt := time.Now()
	res, err := repo.db.Exec(ctx, `UPDATE hex_fwk.product SET name = $2, short_description = $3, description = $4, 
	price = $5, updated_at = $6, quantity = $7, category_id = $8 WHERE id = $1`,
		id, product.Name, product.ShortDescription, product.Description, product.Price, updatedAt, product.Quantity, product.Category.Id)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}
