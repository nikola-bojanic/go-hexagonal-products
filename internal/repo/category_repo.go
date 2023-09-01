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

var _ ports.CategoryRepo = (*CategoryRepository)(nil)

type CategoryRepository struct {
	db *database.DB
}

func NewCategoryRepository(db *database.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (repo *CategoryRepository) GetAllCategories(ctx context.Context) (*[]domain.Category, error) {
	var categories []domain.Category
	rows, err := repo.db.Query(ctx, `SELECT * FROM hex_fwk.category`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var category domain.Category
		err = rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return &categories, nil
}

func (repo *CategoryRepository) FindCategoryById(ctx context.Context, id int64) (*domain.Category, error) {
	var category domain.Category

	// err := repo.db.QueryRow(ctx, `SELECT category_id, category_name, created_at, updated_at FROM hex_fwk.category WHERE category_id = $1`, id).
	// 	StructScan(&category)

	err := repo.db.QueryRow(ctx, `SELECT category_id, category_name, created_at, updated_at FROM hex_fwk.category WHERE category_id = $1`, id).
		Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("Category not found")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (repo *CategoryRepository) InsertCategory(ctx context.Context, category *domain.Category) (int64, error) {
	var id int64
	err := repo.db.QueryRow(ctx, `INSERT INTO hex_fwk.category (category_name) VALUES ($1) RETURNING category_id`,
		category.Name).
		Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *CategoryRepository) DeleteCategory(ctx context.Context, id int64) (int64, error) {
	res, err := repo.db.Exec(ctx, `DELETE FROM hex_fwk.category WHERE category_id = $1`,
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

func (repo *CategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category, id int64) (int64, error) {
	updatedAt := time.Now()
	res, err := repo.db.Exec(ctx, `UPDATE hex_fwk.category SET category_name = $2, updated_at = $3 WHERE category_id = $1`,
		id, category.Name, updatedAt)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}
