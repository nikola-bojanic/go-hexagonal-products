package usecases

import (
	"context"

	domain "github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	ports "github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/pkg/errors"
)

// Check this service satisfies interface
var _ ports.UserUsecase = (*UserService)(nil)

type UserService struct {
	userRepo *repo.UserRepository
}

func NewUserService(userRepo *repo.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, user *domain.User) error {
	err := s.userRepo.Insert(ctx, user)
	if err != nil {
		return errors.Wrap(err, "Failed to register user")
	}
	return nil
}

func (s *UserService) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	return user, nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, egwUser *domain.User) error {
	err := s.userRepo.Update(ctx, egwUser)
	if err != nil {
		return err
	}
	return nil
}
