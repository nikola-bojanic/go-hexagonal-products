package user

import (
	"time"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
)

type UserModel struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	PasswordHash string `json:"password_hash"`

	CreatedAt time.Time `json:"created_at"`
}

func (e *UserModel) FromDomain(user *domain.User) {
	// sanity checks
	if e == nil || user == nil {
		return
	}

	e.ID = user.ID
	e.Email = user.Email
	e.Name = user.Name
	e.Surname = user.Surname
	e.CreatedAt = user.CreatedAt
	// do not populate the password hash, because we do not wish to expose that when loading from the domain
}

func (e *UserModel) ToDomain() *domain.User {
	if e == nil {
		return &domain.User{}
	}

	return &domain.User{
		ID:           e.ID,
		Email:        e.Email,
		CreatedAt:    e.CreatedAt,
		Name:         e.Name,
		Surname:      e.Surname,
		PasswordHash: e.PasswordHash,
	}
}
