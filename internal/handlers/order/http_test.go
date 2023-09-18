package order

// import (
// 	"net/http"
// 	"testing"

// 	"github.com/emicklei/go-restful/v3"
// 	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/app"
// 	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/usecases"
// 	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
// 	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/testutil"
// 	"github.com/stretchr/testify/suite"
// )

// var testApp *app.App

// type HttpSuite struct {
// 	suite.Suite
// 	OrderHttpSvc OrderHttpHandler
// 	wsContainer  *restful.Container
// }

// func (suite *HttpSuite) SetupTest() {

// }

// func (suite *HttpSuite) TearDownTest() {
// 	testutil.CleanUpTables(*testApp.DB)
// }

// func (suite *HttpSuite) SetupSuite() {

// 	testApp = testutil.InitTestApp()
// 	suite.wsContainer = restful.NewContainer()
// 	http.Handle("/", suite.wsContainer)
// 	realUserRep := repo.NewUserRepository(testApp.DB)
// 	realUserSvc := usecases.NewUserService(realUserRep)
// 	realCategoryRep := repo.NewCategoryRepository(testApp.DB)
// 	realCategorySvc := usecases.NewCategoryService(realCategoryRep)
// 	realProductRep := repo.NewProductRepository(testApp.DB)
// 	realProductSvc := usecases.NewProductService(realProductRep)
// 	realOrderRep := repo.NewOrderRepository(testApp.DB)
// 	realOrderSvc := usecases.NewOrderService(realOrderRep, realProductRep)
// 	suite.OrderHttpSvc = *NewOrderHandler(realOrderSvc, realProductSvc, realCategorySvc, realUserSvc, suite.wsContainer)
// }

// func TestOrderTestSuite(t *testing.T) {
// 	suite.Run(t, new(HttpSuite))
// }

// func (suite *HttpSuite) TestCreateOrder() {
// 	categoryName := "test"
// 	cId, err := suite.OrderHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
// 		Name: categoryName,
// 	})
// 	if err != nil {
// 		suite.T().Fatalf("Error creating test category: %s", err)
// 	}
// 	productName := "test"
// 	productShortDescription := "t"
// 	productDescription := "testing"
// 	productPrice := float32(100.0)
// 	productQuantity := 100
// 	productCategory := &domain.Category{Id: int(cId)}
// 	pId, err := suite.OrderHttpSvc.productSvc.CreateProduct(context.TODO(), &domain.Product{
// 		Name:             productName,
// 		ShortDescription: productShortDescription,
// 		Description:      productDescription,
// 		Price:            productPrice,
// 		Quantity:         productQuantity,
// 		Category:         productCategory,
// 	})
// 	if err != nil {
// 		suite.T().Fatalf("Error creating test product: %s", err)
// 	}
// 	orderStatus := ""
// 	orderProducts := &[]OrderedProductModel{{ProductId: pId, Quantity: 10}}
// 	responseRec := testutil.MakeRequest(suite.wsContainer, "POST", "/order", OrderRequest{Status: orderStatus, Products: orderProducts}, nil)
// 	var response *domain.Order
// 	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
// 	if err != nil {
// 		suite.T().Fatalf("Error unmarshalling order response: %s", err)
// 	}
// 	var products []domain.OrderedProduct
// 	for _, product := range *orderProducts {
// 		odreredProduct := product.ToDomain()
// 		products = append(products, *odreredProduct)
// 	}
// 	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
// 	assert.Equal(suite.T(), &products, response.ProductItems)
// 	assert.Equal(suite.T(), "CREATED", response.Status)
// }

// func (suite *HttpSuite) TestUpdateOrderStatus() {
// 	categoryName := "test"
// 	cId, err := suite.OrderHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
// 		Name: categoryName,
// 	})
// 	if err != nil {
// 		suite.T().Fatalf("Error creating test category: %s", err)
// 	}
// 	productName := "test"
// 	productShortDescription := "t"
// 	productDescription := "testing"
// 	productPrice := float32(100.0)
// 	productQuantity := 100
// 	productCategory := &domain.Category{Id: int(cId)}
// 	pId, err := suite.OrderHttpSvc.productSvc.CreateProduct(context.TODO(), &domain.Product{
// 		Name:             productName,
// 		ShortDescription: productShortDescription,
// 		Description:      productDescription,
// 		Price:            productPrice,
// 		Quantity:         productQuantity,
// 		Category:         productCategory,
// 	})
// 	if err != nil {
// 		suite.T().Fatalf("Error creating test product: %s", err)
// 	}
// 	orderStatus := ""
// 	orderProducts := &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}}
// 	created, err := suite.OrderHttpSvc.orderSvc.CreateOrder(context.TODO(), &domain.Order{
// 		Status:       orderStatus,
// 		ProductItems: orderProducts,
// 	})
// 	if err != nil {
// 		suite.T().Fatalf("Error creating test order: %s", err)
// 	}
// 	updateID := created.ID
// 	updateStatus := "PENDING"
// 	updateProducts := &[]OrderedProductModel{{ProductId: pId, Quantity: 10}}

// 	responseRec := testutil.MakeRequest(suite.wsContainer, "PUT", "/order/status", OrderRequest{
// 		ID:       updateID,
// 		Status:   updateStatus,
// 		Products: updateProducts,
// 	}, nil)
// 	var response *domain.Order
// 	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
// 	if err != nil {
// 		suite.T().Fatalf("Error unmarshalling order response: %s", err)
// 	}
// 	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
// 	assert.Equal(suite.T(), updateStatus, response.Status)
// }

// func (suite *HttpSuite) TestDeleteOrder() {
// 	userEmail := "testy@email.com"
// 	// register the user
// 	userPass := "password123"
// 	passHash, _ := bcrypt.GenerateFromPassword([]byte(userPass), 10)
// 	suite.OrderHttpSvc.userSvc.RegisterUser(context.TODO(), &domain.User{
// 		Email:        userEmail,
// 		PasswordHash: string(passHash),
// 	})
// 	categoryName := "test"
// 	cId, err := suite.OrderHttpSvc.categorySvc.CreateCategory(context.TODO(), &domain.Category{
// 		Name: categoryName,
// 	})
// 	if err != nil {
// 		suite.T().Fatalf("Error creating test category: %s", err)
// 	}
// 	productName := "test"
// 	productShortDescription := "t"
// 	productDescription := "testing"
// 	productPrice := float32(100.0)
// 	productQuantity := 100
// 	productCategory := &domain.Category{Id: int(cId)}
// 	pId, err := suite.OrderHttpSvc.productSvc.CreateProduct(context.TODO(), &domain.Product{
// 		Name:             productName,
// 		ShortDescription: productShortDescription,
// 		Description:      productDescription,
// 		Price:            productPrice,
// 		Quantity:         productQuantity,
// 		Category:         productCategory,
// 	})
// 	if err != nil {
// 		suite.T().Fatalf("Error creating test product: %s", err)
// 	}
// 	user, err := suite.OrderHttpSvc.userSvc.FindByEmail(context.TODO(), userEmail)
// 	if err != nil {
// 		suite.T().Fatalf("cannot fetch user: %s", err)
// 	}
// 	orderStatus := ""
// 	orderProducts := &[]domain.OrderedProduct{{ProductId: pId, Quantity: 10}}
// 	created, err := suite.OrderHttpSvc.orderSvc.CreateOrder(context.TODO(), &domain.Order{
// 		Status:       orderStatus,
// 		ProductItems: orderProducts,
// 		User:         user,
// 	})
// 	if err != nil {
// 		suite.T().Fatalf("Error creating test order: %s", err)
// 	}
// 	deleteID := created.ID
// 	deleteStatus := "CREATED"
// 	deleteProducts := &[]OrderedProductModel{{ProductId: pId, Quantity: 10}}

// 	responseRec := testutil.MakeRequest(suite.wsContainer, "DELETE", "/order", OrderRequest{
// 		ID:       deleteID,
// 		Status:   deleteStatus,
// 		Products: deleteProducts,
// 		UserId:   user.ID,
// 	}, nil)
// 	var response *domain.Order
// 	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
// 	if err != nil {
// 		suite.T().Fatalf("Error unmarshalling order response: %s", err)
// 	}
// 	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
// 	assert.Equal(suite.T(), deleteID, response.ID)

// }
