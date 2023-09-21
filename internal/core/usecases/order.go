package usecases

import (
	"context"
	"strconv"

	"github.com/jung-kurt/gofpdf"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/pkg/errors"
)

var _ ports.OrderUsecase = (*OrderService)(nil)

type OrderService struct {
	orderRepo   *repo.OrderRepository
	productRepo *repo.ProductRepository
	userRepo    *repo.UserRepository
}

func NewOrderService(orderRepo *repo.OrderRepository, productRepo *repo.ProductRepository, userRepo *repo.UserRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

func (s *OrderService) FindOrderById(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.orderRepo.FindOrderById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve an product")
	}
	return order, nil
}
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	order.Status = "CREATED"
	for _, item := range *order.ProductItems {
		product, err := s.productRepo.FindProductById(ctx, item.ProductId)
		if err != nil {
			return nil, errors.New("product doesn't exist")
		}
		if item.Quantity <= 0 {
			return nil, errors.New("invalid order")
		}
		if product.Quantity < item.Quantity {
			return nil, errors.New("not enough items")
		}
		product.Quantity = product.Quantity - item.Quantity
		_, err = s.productRepo.UpdateProduct(ctx, product, item.ProductId)
		if err != nil {
			return nil, errors.New("error updating product quantity")
		}
	}
	created, err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create an order")
	}
	return created, nil
}
func (s *OrderService) UpdateOrderStatus(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	validStatus := order.Status == "" || order.Status == "CREATED" || order.Status == "PENDING" || order.Status == "COMPLETED" || order.Status == "CLOSED"
	if !validStatus {
		return nil, errors.New("invalid order status")
	}
	for _, item := range *order.ProductItems {
		product, err := s.productRepo.FindProductById(ctx, item.ProductId)
		if err != nil {
			return nil, errors.New("product doesn't exist")
		}
		if item.Quantity <= 0 {
			return nil, errors.New("invalid order")
		}
		if product.Quantity < item.Quantity {
			return nil, errors.New("not enough items")
		}
		product.Quantity = product.Quantity - item.Quantity
		_, err = s.productRepo.UpdateProduct(ctx, product, item.ProductId)
		if err != nil {
			return nil, errors.New("error updating product quantity")
		}
	}
	updated, err := s.orderRepo.UpdateOrderStatus(ctx, order)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update an order")
	}
	return updated, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, order *domain.Order) error {
	err := s.orderRepo.DeleteOrder(ctx, order)
	if err != nil {
		return errors.Wrap(err, "failed to update an order")
	}
	return nil
}

func (s *OrderService) GeneratePdf(ctx context.Context, order *domain.Order) error {
	pdf := s.newReport(ctx, order)
	pdf = header(pdf, []string{"Product No.", "Name", "Quantity", "Unit Price", "Total Price"})

	pdf = s.tableContent(ctx, pdf, order)
	if pdf.Err() {
		return pdf.Error()
	}
	err := pdf.OutputFileAndClose("pdf/order/order_[" + order.ID + "].pdf")

	if err != nil {
		return err
	}
	return nil
}
func (s *OrderService) tableContent(ctx context.Context, pdf *gofpdf.Fpdf, order *domain.Order) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 16)
	pdf.SetFillColor(255, 255, 255)

	for i, orderedProduct := range *order.ProductItems {
		product, _ := s.productRepo.FindProductById(ctx, orderedProduct.ProductId)

		pdf.CellFormat(40, 10, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, product.Name, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, strconv.Itoa(orderedProduct.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, strconv.FormatFloat(float64(product.Price), 'f', 2, 64), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, strconv.FormatFloat(float64(orderedProduct.Quantity)*float64(product.Price), 'f', 2, 64), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(-1)

	return pdf
}

func (s *OrderService) newReport(ctx context.Context, order *domain.Order) *gofpdf.Fpdf {
	user, _ := s.userRepo.FindByID(ctx, order.User.ID)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Times", "B", 20)
	pdf.Cell(40, 10, "Order ID: "+order.ID)
	pdf.Ln(15)

	pdf.SetFont("Times", "B", 20)
	pdf.Cell(40, 10, "User ID: "+user.ID)
	pdf.Ln(-1)

	pdf.SetFont("Times", "B", 20)
	pdf.Cell(40, 10, "User name: "+user.Name+" "+user.Surname)
	pdf.Ln(-1)

	pdf.SetFont("Times", "B", 20)
	pdf.Cell(40, 10, "User e-mail: "+user.Email)
	pdf.Ln(20)
	return pdf
}

func header(pdf *gofpdf.Fpdf, headerText []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 16)
	pdf.SetFillColor(240, 240, 240)

	for _, str := range headerText {
		pdf.CellFormat(40, 10, str, "1", 0, "C", true, 0, "")
	}

	pdf.Ln(-1)

	return pdf
}
