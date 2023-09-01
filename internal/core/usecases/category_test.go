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

type CategorySuite struct {
	suite.Suite
	categoryRep *repo.CategoryRepository
	categorySvc *CategoryService
}

func (suite *CategorySuite) SetupTest() {

}

func (suite *CategorySuite) TearDownTest() {

}

func (suite *CategorySuite) SetupSuite() {
	app := testutil.InitTestApp()
	suite.categoryRep = repo.NewCategoryRepository(app.DB)
	suite.categorySvc = NewCategoryService(suite.categoryRep)
}

func TestCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}

func (suite *CategorySuite) TestNamelessCategoryCreation() {
	var category domain.Category
	_, err := suite.categorySvc.CreateCategory(context.TODO(), &category)
	assert.NotNil(suite.T(), err)
}
