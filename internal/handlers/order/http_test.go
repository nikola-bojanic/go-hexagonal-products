package order

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/app"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/usecases"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testApp *app.App

type HttpSuite struct {
	suite.Suite
	OrderHttpSvc OrderHttpHandler
	CategoryRepo *repo.CategoryRepository
	ProductRepo  *repo.ProductRepository
	OrderRepo    *repo.OrderRepository
	wsContainer  *restful.Container
}

func (suite *HttpSuite) SetupTest() {
}

func (suite *HttpSuite) TearDownTest() {
	testutil.CleanUpTables(*testApp.DB)
}

func (suite *HttpSuite) SetupSuite() {

	testApp = testutil.InitTestApp()
	suite.wsContainer = restful.NewContainer()
	http.Handle("/", suite.wsContainer)
	realProductRep := repo.NewProductRepository(testApp.DB)
	realOrderRep := repo.NewOrderRepository(testApp.DB)
	realOrderSvc := usecases.NewOrderService(realOrderRep, realProductRep)
	suite.OrderRepo = realOrderRep
	suite.ProductRepo = realProductRep
	suite.CategoryRepo = repo.NewCategoryRepository(testApp.DB)
	suite.OrderHttpSvc = *NewOrderHandler(realOrderSvc, suite.wsContainer)
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (suite *HttpSuite) TestCreateOrder() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.CategoryRepo.InsertCategory(context.TODO(), &testCategory)
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
	pId, err := suite.ProductRepo.InsertProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test order: %s", err)
	}
	testOrder := OrderRequest{
		Status:   "",
		Products: &[]OrderedProductModel{{ProductId: pId, Quantity: 10}},
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "POST", "/order", testOrder, nil)
	var response *domain.Order
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling order response: %s", err)
	}
	var products []domain.OrderedProduct
	for _, product := range *testOrder.Products {
		odreredProduct := product.ToDomain()
		products = append(products, *odreredProduct)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.Equal(suite.T(), &products, response.ProductItems)
	assert.Equal(suite.T(), "CREATED", response.Status)
	suite.OrderRepo.DeleteOrder(context.TODO(), response)
	suite.ProductRepo.DeleteProduct(context.TODO(), pId)
	suite.CategoryRepo.DeleteCategory(context.TODO(), cId)
}

func (suite *HttpSuite) TestUpdateOrderStatus() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.CategoryRepo.InsertCategory(context.TODO(), &testCategory)
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
	pId, err := suite.ProductRepo.InsertProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test order: %s", err)
	}
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}},
	}
	created, err := suite.OrderRepo.CreateOrder(context.TODO(), &testOrder)
	if err != nil {
		suite.T().Fatalf("Error creating test order: %s", err)
	}
	updateOrder := OrderRequest{
		ID:       created.ID,
		Status:   "PENDING",
		Products: &[]OrderedProductModel{{ProductId: pId, Quantity: 10}},
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "PUT", "/order/status", updateOrder, nil)
	var response *domain.Order
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling order response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.Equal(suite.T(), updateOrder.Status, response.Status)
	suite.OrderRepo.DeleteOrder(context.TODO(), response)
	suite.ProductRepo.DeleteProduct(context.TODO(), pId)
	suite.CategoryRepo.DeleteCategory(context.TODO(), cId)
}

func (suite *HttpSuite) TestDeleteOrder() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.CategoryRepo.InsertCategory(context.TODO(), &testCategory)
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
	pId, err := suite.ProductRepo.InsertProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	testOrder := domain.Order{
		Status:       "",
		ProductItems: &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}},
	}
	created, err := suite.OrderRepo.CreateOrder(context.TODO(), &testOrder)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	deleteOrder := OrderRequest{
		ID:       created.ID,
		Status:   "CREATED",
		Products: &[]OrderedProductModel{{ProductId: pId, Quantity: 10}},
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "DELETE", "/order", deleteOrder, nil)
	var response *domain.Order
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling order response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.Equal(suite.T(), deleteOrder.ID, response.ID)
	suite.ProductRepo.DeleteProduct(context.TODO(), pId)
	suite.CategoryRepo.DeleteCategory(context.TODO(), cId)
}
