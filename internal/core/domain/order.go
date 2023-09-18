package domain

import (
	"fmt"
	"time"
)

type Order struct {
	ID           string            `json:"id"`
	ProductItems *[]OrderedProduct `json:"product_items"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	User         *User             `json:"user"`
}
type OrderedProduct struct {
	ProductId int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

func (e *Order) ToString() string {
	return fmt.Sprintf("%s %v %s %v", e.ID, e.ProductItems, e.Status, e.User)
}
