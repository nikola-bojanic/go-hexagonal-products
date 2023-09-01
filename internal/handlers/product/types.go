package product

import "github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/handlers/category"

type Response struct {
	ID      int64
	Message string
}

type ProductRequest struct {
	Name             string
	ShortDescription string
	Description      string
	Price            float32
	Quantity         int
	Category         *category.CategoryModel
}
