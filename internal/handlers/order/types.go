package order

type Response struct {
	ID   int64
	Name string
}

type OrderResponse struct {
	Status   string
	Products *[]OrderedProduct
}

type OrderedProduct struct {
	ID       string
	Name     string
	Quantity int
}

type OrderRequest struct {
	OrderId  string
	Status   string
	Products *[]OrderedProductModel
}
