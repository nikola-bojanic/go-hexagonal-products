package domain

import (
	"fmt"
	"time"
)

type Product struct {
	ProductId        int       `json:"productId"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"shortDescription"`
	Description      string    `json:"description"`
	Price            float32   `json:"price"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Quantity         int       `json:"quantity"`
	Category         *Category `json:"category"`
}

func (e *Product) ToString() string {
	return fmt.Sprintf("%d %s %s %s %f %d", e.ProductId, e.Name, e.ShortDescription, e.Description, e.Price, e.Quantity)
}
