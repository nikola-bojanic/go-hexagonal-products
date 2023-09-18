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
	categoryName := "test"
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
		Name: categoryName,
	})
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	productName := "test"
	productShortDescription := "t"
	productDescription := "testing"
	productPrice := float32(100.0)
	productQuantity := 1
	productCategory := &domain.Category{Id: int(cId)}
	_, err = suite.productHttpSvc.productSvc.CreateProduct(context.TODO(), &domain.Product{
		Name:             productName,
		ShortDescription: productShortDescription,
		Description:      productDescription,
		Price:            productPrice,
		Quantity:         productQuantity,
		Category:         productCategory,
	})
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
}
func (suite *HttpSuite) TestGetProduct() {
	categoryName := "test"
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
		Name: categoryName,
	})
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	productName := "test"
	productShortDescription := "t"
	productDescription := "testing"
	productPrice := float32(100.0)
	productQuantity := 1
	productCategory := &domain.Category{Id: int(cId)}
	pId, err := suite.productHttpSvc.productSvc.CreateProduct(context.TODO(), &domain.Product{
		Name:             productName,
		ShortDescription: productShortDescription,
		Description:      productDescription,
		Price:            productPrice,
		Quantity:         productQuantity,
		Category:         productCategory,
	})
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
	assert.Equal(suite.T(), productName, response.Name)
	assert.Equal(suite.T(), productCategory.Id, response.Category.Id)
	assert.Equal(suite.T(), productDescription, response.Description)
	assert.Equal(suite.T(), productPrice, response.Price)
	assert.Equal(suite.T(), productShortDescription, response.ShortDescription)
	assert.Equal(suite.T(), productQuantity, response.Quantity)
}
func (suite *HttpSuite) TestCreateProduct() {
	categoryName := "test"
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
		Name: categoryName,
	})
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	productName := "test"
	productShortDescription := "t"
	productDescription := "testing"
	productPrice := float32(100.0)
	productQuantity := 1
	productCategory := &domain.Category{Id: int(cId)}
	responseRec := testutil.MakeRequest(suite.wsContainer, "POST", "/product", &domain.Product{
		Name:             productName,
		ShortDescription: productShortDescription,
		Description:      productDescription,
		Price:            productPrice,
		Quantity:         productQuantity,
		Category:         productCategory,
	}, nil)
	var response Response
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling category response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.Equal(suite.T(), "product created", response.Message)
}
func (suite *HttpSuite) TestUpdateProduct() {
	categoryName := "test"
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
		Name: categoryName,
	})
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	productName := "test"
	productShortDescription := "t"
	productDescription := "testing"
	productPrice := float32(100.0)
	productQuantity := 1
	productCategory := &domain.Category{Id: int(cId)}
	pId, err := suite.productHttpSvc.productSvc.CreateProduct(context.TODO(), &domain.Product{
		Name:             productName,
		ShortDescription: productShortDescription,
		Description:      productDescription,
		Price:            productPrice,
		Quantity:         productQuantity,
		Category:         productCategory,
	})
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	updateCname := "updated"
	uCId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
		Name: updateCname,
	})
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	updateName := "test2"
	updateShortDescription := "t2"
	updateDescription := "testing2"
	updatePrice := float32(200.0)
	updateQuantity := 2
	updateCategory := &domain.Category{Id: int(uCId)}

	responseRec := testutil.MakeRequest(suite.wsContainer, "PUT", "/product/"+strconv.Itoa(int(pId)), domain.Product{
		Name:             updateName,
		ShortDescription: updateShortDescription,
		Description:      updateDescription,
		Price:            updatePrice,
		Quantity:         updateQuantity,
		Category:         updateCategory,
	}, nil)
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
}
func (suite *HttpSuite) TestDeleteProduct() {
	categoryName := "test"
	cId, err := suite.productHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
		Name: categoryName,
	})
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	productName := "test"
	productShortDescription := "t"
	productDescription := "testing"
	productPrice := float32(100.0)
	productQuantity := 1
	productCategory := &domain.Category{Id: int(cId)}
	pId, err := suite.productHttpSvc.productSvc.CreateProduct(context.TODO(), &domain.Product{
		Name:             productName,
		ShortDescription: productShortDescription,
		Description:      productDescription,
		Price:            productPrice,
		Quantity:         productQuantity,
		Category:         productCategory,
	})
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
}
