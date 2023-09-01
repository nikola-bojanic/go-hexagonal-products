package repo

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/pkg/errors"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/database"
)

var ErrDuplicateEmail = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")

// Verify the impl matches the interface
var _ ports.UserRepo = (*UserRepository)(nil)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Insert(ctx context.Context, user *domain.User) error {
	_, err := repo.db.Exec(ctx,
		"INSERT INTO hex_fwk.user (email, first_name, surname,password_hash) VALUES ($1, $2, $3, $4)",
		user.Email, user.Name, user.Surname, user.PasswordHash)
	if err != nil {
		alreadyExists, _ := regexp.Match(`user_email_key`, []byte(err.Error()))
		if alreadyExists {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (repo *UserRepository) Update(ctx context.Context, user *domain.User) error {
	// update, and reflect changes in the struct
	err := repo.db.QueryRow(ctx,
		`UPDATE hex_fwk.user SET
			first_name = $1,
			surname = $2
		 WHERE id = $3
		 RETURNING id, first_name, surname, email`,
		user.Name, user.Surname, user.ID).StructScan(user)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	err := repo.db.
		QueryRow(ctx, `SELECT id, email, first_name, surname, password_hash FROM hex_fwk.user WHERE id = $1`, id).
		StructScan(&user)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	err := repo.db.QueryRow(ctx,
		`SELECT id, email, first_name, surname, password_hash FROM hex_fwk.user WHERE email = $1`,
		email).
		StructScan(&user)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
