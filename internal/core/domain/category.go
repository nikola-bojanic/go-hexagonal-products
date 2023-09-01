package domain

import (
	"fmt"
	"time"
)

type Category struct {
	Id        int       `json:"categoryId"`
	Name      string    `json:"categoryName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (e *Category) ToString() string {
	return fmt.Sprintf("%d %s", e.Id, e.Name)
}
