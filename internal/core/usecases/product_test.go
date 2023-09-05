package usecases

import (
	"context"
	"testing"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductSuite struct {
	suite.Suite
	productRep  *repo.ProductRepository
	productSvc  *ProductService
	categoryRep *repo.CategoryRepository
	categorySvc *CategoryService
}

func (suite *ProductSuite) SetupTest() {

}

func (suite *ProductSuite) TearDownTest() {

}

func (suite *ProductSuite) SetupSuite() {
	app := testutil.InitTestApp()
	suite.productRep = repo.NewProductRepository(app.DB)
	suite.productSvc = NewProductService(suite.productRep)
	suite.categoryRep = repo.NewCategoryRepository(app.DB)
	suite.categorySvc = NewCategoryService(suite.categoryRep)
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductSuite))
}

func (suite *ProductSuite) TestCreateProductWithoutCategory() {
	product := domain.Product{
		Name:             "",
		ShortDescription: "",
		Description:      "",
		Price:            0,
		Quantity:         0,
		Category:         &domain.Category{},
	}
	_, err := suite.productSvc.CreateProduct(context.TODO(), &product)
	assert.NotNil(suite.T(), err)
}
func (suite *ProductSuite) TestCreateProductWithInvalidCategory() {
	product := domain.Product{
		Name:             "",
		ShortDescription: "",
		Description:      "",
		Price:            0,
		Quantity:         0,
		Category:         &domain.Category{Id: 500},
	}
	_, err := suite.productSvc.CreateProduct(context.TODO(), &product)
	assert.NotNil(suite.T(), err)
}

func (suite *ProductSuite) TestGetProducts() {
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
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productRep.InsertProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	products, err := suite.productSvc.GetAllProducts(context.TODO())
	if err != nil {
		suite.T().Fatal(err)
	}
	assert.NotNil(suite.T(), products)
	suite.categorySvc.DeleteCategory(context.TODO(), cId)
	suite.productRep.DeleteProduct(context.TODO(), pId)
}
func (suite *ProductSuite) TestGetProduct() {
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
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productRep.InsertProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	product, err := suite.productSvc.FindProductById(context.TODO(), pId)
	if err != nil {
		suite.T().Fatal(err)
	}
	assert.Equal(suite.T(), testProduct.Name, product.Name)
	assert.Equal(suite.T(), testProduct.ShortDescription, product.ShortDescription)
	assert.Equal(suite.T(), testProduct.Description, product.Description)
	assert.Equal(suite.T(), testProduct.Price, product.Price)
	assert.Equal(suite.T(), testProduct.Quantity, product.Quantity)
	assert.Equal(suite.T(), testProduct.Category.Id, product.Category.Id)

	suite.categorySvc.DeleteCategory(context.TODO(), cId)
	suite.productRep.DeleteProduct(context.TODO(), pId)
}
func (suite *ProductSuite) TestCreateProduct() {
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
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productSvc.CreateProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatal(err)
	}
	zeroId := int64(0)
	assert.NotEqual(suite.T(), zeroId, pId)
	suite.categorySvc.DeleteCategory(context.TODO(), cId)
	suite.productRep.DeleteProduct(context.TODO(), pId)
}
func (suite *ProductSuite) TestUpdateProduct() {
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
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productRep.InsertProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	updateProduct := domain.Product{
		Name:             "test2",
		ShortDescription: "t2",
		Description:      "testing2",
		Price:            200.0,
		Quantity:         2,
		Category:         &domain.Category{Id: int(cId)},
	}
	rows, err := suite.productSvc.UpdateProduct(context.TODO(), &updateProduct, pId)
	if err != nil {
		suite.T().Fatal(err)
	}
	noRows := int64(0)
	assert.NotEqual(suite.T(), noRows, rows)
	suite.categorySvc.DeleteCategory(context.TODO(), cId)
	suite.productRep.DeleteProduct(context.TODO(), pId)
}
func (suite *ProductSuite) TestDeleteProduct() {
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
		Price:            100.0,
		Quantity:         1,
		Category:         &domain.Category{Id: int(cId)},
	}
	pId, err := suite.productRep.InsertProduct(context.TODO(), &testProduct)
	if err != nil {
		suite.T().Fatalf("Error creating test product: %s", err)
	}
	rows, err := suite.productSvc.DeleteProduct(context.TODO(), pId)
	if err != nil {
		suite.T().Fatal(err)
	}
	noRows := int64(0)
	assert.NotEqual(suite.T(), noRows, rows)
	suite.categorySvc.DeleteCategory(context.TODO(), cId)
}
