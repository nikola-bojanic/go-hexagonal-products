package order

type Response struct {
	ID   int64
	Name string
}

type OrderResponse struct {
	Status   string
	UserId   string
	Products *[]OrderedProduct
}

type OrderedProduct struct {
	ID       string
	Name     string
	Quantity int
}

type OrderRequest struct {
	ID       string
	Status   string
	UserId   string
	Products *[]OrderedProductModel `json:"product_items"`
}
