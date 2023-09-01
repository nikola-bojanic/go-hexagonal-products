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
	productRep *repo.ProductRepository
	productSvc *ProductService
}

func (suite *ProductSuite) SetupTest() {

}

func (suite *ProductSuite) TearDownTest() {

}

func (suite *ProductSuite) SetupSuite() {
	app := testutil.InitTestApp()
	suite.productRep = repo.NewProductRepository(app.DB)
	suite.productSvc = NewProductService(suite.productRep)
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
