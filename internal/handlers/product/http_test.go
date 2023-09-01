package product

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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
	productHttpSvc ProductHttpHandler
	wsContainer    *restful.Container
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
	realCategoryRep := repo.NewCategoryRepository(testApp.DB)
	realCategorySvc := usecases.NewCategoryService(realCategoryRep)
	realProductRep := repo.NewProductRepository(testApp.DB)
	realProductSvc := usecases.NewProductService(realProductRep)
	suite.productHttpSvc = *NewProductHandler(realProductSvc, realCategorySvc, suite.wsContainer)
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}

func (suite *HttpSuite) TestGetProducts() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productHttpSvc.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "GET", "/product", nil, nil)
	var response []ProductModel
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.NotNil(suite.T(), response)
	suite.productHttpSvc.productSvc.DeleteProduct(context.TODO(), pId)
	suite.productHttpSvc.categorySvc.DeleteCategory(context.TODO(), cId)
}
func (suite *HttpSuite) TestGetProduct() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productHttpSvc.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "GET", "/product/"+strconv.Itoa(int(pId)), nil, nil)
	var response ProductModel
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.Equal(suite.T(), testProduct.Name, response.Name)
	assert.Equal(suite.T(), testProduct.Category.Id, response.Category.Id)
	assert.Equal(suite.T(), testProduct.Description, response.Description)
	assert.Equal(suite.T(), testProduct.Price, response.Price)
	assert.Equal(suite.T(), testProduct.ShortDescription, response.ShortDescription)
	assert.Equal(suite.T(), testProduct.Quantity, response.Quantity)
	suite.productHttpSvc.productSvc.DeleteProduct(context.TODO(), pId)
	suite.productHttpSvc.categorySvc.DeleteCategory(context.TODO(), cId)
}
func (suite *HttpSuite) TestCreateProduct() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "POST", "/product", testProduct, nil)
	var response Response
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling category response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.Equal(suite.T(), "product created", response.Message)
	suite.productHttpSvc.productSvc.DeleteProduct(context.TODO(), response.ID)
	suite.productHttpSvc.categorySvc.DeleteCategory(context.TODO(), cId)
}
func (suite *HttpSuite) TestUpdateProduct() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productHttpSvc.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	updateCategory := domain.Category{
		Name: "test2",
	}
	uCId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &updateCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	updateProduct := domain.Product{
		Name:             "test2",
		ShortDescription: "t2",
		Description:      "testing2",
		Price:            200.0,
		Quantity:         2,
		Category:         &domain.Category{Id: int(uCId)},
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "PUT", "/product/"+strconv.Itoa(int(pId)), updateProduct, nil)
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	var response Response
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling product response: %s", err)
	}
	message := "product updated"
	assert.Equal(suite.T(), message, response.Message)
	rowsAffected := int64(1)
	assert.Equal(suite.T(), rowsAffected, response.ID)
	suite.productHttpSvc.productSvc.DeleteProduct(context.TODO(), pId)
	suite.productHttpSvc.categorySvc.DeleteCategory(context.TODO(), cId)
}
func (suite *HttpSuite) TestDeleteProduct() {
	testCategory := domain.Category{
		Name: "test",
	}
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	testProduct := domain.Product{
		Name:             "test",
		ShortDescription: "t",
		Description:      "testing",
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productHttpSvc.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "DELETE", "/product/"+strconv.Itoa(int(pId)), nil, nil)
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	message := "product deleted"
	var response Response
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling category response: %s", err)
	}
	assert.Equal(suite.T(), message, response.Message)
	rowsAffected := int64(1)
	assert.Equal(suite.T(), rowsAffected, response.ID)
	suite.productHttpSvc.productSvc.DeleteProduct(context.TODO(), pId)
	suite.productHttpSvc.categorySvc.DeleteCategory(context.TODO(), cId)
}
