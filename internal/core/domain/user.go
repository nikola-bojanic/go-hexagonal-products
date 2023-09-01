package domain

import (
	"fmt"
	"time"
)

type User struct {
	ID           string `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	Name         string `json:"name" db:"first_name"`
	Surname      string `json:"surname" db:"surname"`
	PasswordHash string `json:"password_hash" db:"password_hash"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func NewUser(id string, email string, name string, surname string) *User {
	return &User{
		ID:      id,
		Email:   email,
		Name:    name,
		Surname: surname,
	}
}

func (e *User) ToString() string {
	return fmt.Sprintf("#%s %s %s - %s", e.ID, e.Name, e.Surname, e.Email)
}
