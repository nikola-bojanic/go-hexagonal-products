package usecases

import (
	"context"

	domain "github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/pkg/errors"
)

var _ ports.CategoryUsecase = (*CategoryService)(nil)

type CategoryService struct {
	categoryRepo *repo.CategoryRepository
}

func NewCategoryService(categoryRepo *repo.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) GetAllCategories(ctx context.Context) (*[]domain.Category, error) {
	categories, err := s.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve categories")
	}
	return categories, nil
}
func (s *CategoryService) FindCategoryById(ctx context.Context, id int64) (*domain.Category, error) {
	category, err := s.categoryRepo.FindCategoryById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve a category")
	}
	return category, nil
}
func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (int64, error) {
	if len(category.Name) < 1 {
		return 0, errors.New("category doesn't have a name")
	}
	id, err := s.categoryRepo.InsertCategory(ctx, category)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to create a category")
	}
	return id, nil
}
func (s *CategoryService) DeleteCategory(ctx context.Context, id int64) (int64, error) {
	id, err := s.categoryRepo.DeleteCategory(ctx, id)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to delete a category")
	}
	return id, nil
}
func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category, id int64) (int64, error) {
	id, err := s.categoryRepo.UpdateCategory(ctx, category, id)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to edit a category")
	}
	return id, nil
}
