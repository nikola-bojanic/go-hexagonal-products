package usecases

import (
	"context"
	"testing"

	domain "github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OrderSuite struct {
	suite.Suite
	orderRep    *repo.OrderRepository
	orderSvc    *OrderService
	productRep  *repo.ProductRepository
	productSvc  *ProductService
	categoryRep *repo.CategoryRepository
	categorySvc *CategoryService
}

func (suite *OrderSuite) SetupTest() {

}

func (suite *OrderSuite) TearDownTest() {

}

func (suite *OrderSuite) SetupSuite() {
	app := testutil.InitTestApp()
	suite.orderRep = repo.NewOrderRepository(app.DB)
	suite.productRep = repo.NewProductRepository(app.DB)
	suite.orderSvc = NewOrderService(suite.orderRep, suite.productRep)
	suite.productSvc = NewProductService(suite.productRep)
	suite.categoryRep = repo.NewCategoryRepository(app.DB)
	suite.categorySvc = NewCategoryService(suite.categoryRep)
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

func (suite *OrderSuite) TestFindOrderById() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            1000.0,
		Quantity:         100,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}},
	}
	created, err := suite.orderRep.CreateOrder(context.TODO(), &testOrder)
	if err != nil {
		suite.T().Fatalf("Error creating test order: %s", err)
	}
	order, err := suite.orderSvc.FindOrderById(context.TODO(), created.ID)
	if err != nil {
		suite.T().Fatal(err)
	}
	assert.Equal(suite.T(), testOrder.ProductItems, order.ProductItems)
	assert.Equal(suite.T(), "CREATED", order.Status)
	suite.orderRep.DeleteOrder(context.TODO(), order)
	suite.productRep.DeleteProduct(context.TODO(), pId)
	suite.categoryRep.DeleteCategory(context.TODO(), cId)
}

func (suite *OrderSuite) TestCreateOrder() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            1000.0,
		Quantity:         100,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}},
	}
	created, err := suite.orderSvc.CreateOrder(context.TODO(), &testOrder)
	if err != nil {
		suite.T().Fatal(err)
	}
	assert.Equal(suite.T(), testOrder.ProductItems, created.ProductItems)
	assert.Equal(suite.T(), "CREATED", created.Status)
	suite.orderRep.DeleteOrder(context.TODO(), created)
	suite.productRep.DeleteProduct(context.TODO(), pId)
	suite.categoryRep.DeleteCategory(context.TODO(), cId)
}

func (suite *OrderSuite) TestUpdateOrderStatus() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            1000.0,
		Quantity:         100,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}},
	}
	created, err := suite.orderRep.CreateOrder(context.TODO(), &testOrder)
	if err != nil {
		suite.T().Fatalf("Error creating test order: %s", err)
	}
	created.Status = "PENDING"
	created, err = suite.orderRep.UpdateOrderStatus(context.TODO(), created)
	if err != nil {
		suite.T().Fatal(err)
	}
	assert.Equal(suite.T(), "PENDING", created.Status)
	suite.orderRep.DeleteOrder(context.TODO(), created)
	suite.productRep.DeleteProduct(context.TODO(), pId)
	suite.categoryRep.DeleteCategory(context.TODO(), cId)
}

func (suite *OrderSuite) TestDeleteOrder() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            1000.0,
		Quantity:         100,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	testOrder := domain.Order{
		Status:       "CREATED",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}},
	}
	created, err := suite.orderRep.CreateOrder(context.TODO(), &testOrder)
	if err != nil {
		suite.T().Fatalf("Error creating test order: %s", err)
	}
	err = suite.orderSvc.DeleteOrder(context.TODO(), created)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.productRep.DeleteProduct(context.TODO(), pId)
	suite.categoryRep.DeleteCategory(context.TODO(), cId)
}

func (suite *OrderSuite) TestCreateOrderWithInvalidProduct() {
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: int64(555), Quantity: 10}},
	}
	_, err := suite.orderSvc.CreateOrder(context.TODO(), &testOrder)
	assert.NotNil(suite.T(), err)
}
func (suite *OrderSuite) TestCreateOrderWithZeroProductQuantity() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            1000.0,
		Quantity:         100,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 0}},
	}
	_, err = suite.orderSvc.CreateOrder(context.TODO(), &testOrder)
	assert.NotNil(suite.T(), err)
	suite.productRep.DeleteProduct(context.TODO(), pId)
	suite.categoryRep.DeleteCategory(context.TODO(), cId)
}
func (suite *OrderSuite) TestCreateOrderWithInvalidProductQuantity() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            1000.0,
		Quantity:         100,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 101}},
	}
	_, err = suite.orderSvc.CreateOrder(context.TODO(), &testOrder)
	assert.NotNil(suite.T(), err)
	suite.productRep.DeleteProduct(context.TODO(), pId)
	suite.categoryRep.DeleteCategory(context.TODO(), cId)
}
func (suite *OrderSuite) TestInvalidProductStatusUpdate() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            1000.0,
		Quantity:         100,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}},
	}
	created, err := suite.orderRep.CreateOrder(context.TODO(), &testOrder)
	if err != nil {
		suite.T().Fatalf("Error creating test order: %s", err)
	}
	created.Status = "invalid"
	_, err = suite.orderSvc.UpdateOrderStatus(context.TODO(), created)
	assert.NotNil(suite.T(), err)
	suite.orderRep.DeleteOrder(context.TODO(), created)
	suite.productRep.DeleteProduct(context.TODO(), pId)
	suite.categoryRep.DeleteCategory(context.TODO(), cId)
}
