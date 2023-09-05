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

func (suite *CategorySuite) TestGetCategories() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("error creating test category %v", err)
	}
	categories, err := suite.categorySvc.GetAllCategories(context.TODO())
	if err != nil {
		suite.T().Fatal(err)
	}
	assert.NotNil(suite.T(), categories)
	suite.categorySvc.DeleteCategory(context.TODO(), id)
}
func (suite *CategorySuite) TestGetCategory() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("error creating test category %v", err)
	}
	category, err := suite.categorySvc.FindCategoryById(context.TODO(), id)
	if err != nil {
		suite.T().Fatal(err)
	}
	assert.Equal(suite.T(), testCategory.Name, category.Name)
	suite.categorySvc.DeleteCategory(context.TODO(), id)
}
func (suite *CategorySuite) TestCreateCategory() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatal(err)
	}
	zeroId := int64(0)
	assert.NotEqual(suite.T(), zeroId, id)
	suite.categorySvc.DeleteCategory(context.TODO(), id)
}
func (suite *CategorySuite) TestUpdateCategory() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("error creating test category %v", err)
	}
	updateCategory := domain.Category{
		Name: "update",
	}
	rows, err := suite.categorySvc.UpdateCategory(context.TODO(), &updateCategory, id)
	if err != nil {
		suite.T().Fatal(err)
	}
	zeroRows := int64(0)
	assert.NotEqual(suite.T(), zeroRows, rows)
	suite.categorySvc.DeleteCategory(context.TODO(), id)
}
func (suite *CategorySuite) TestDeleteCategory() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("error creating test category %v", err)
	}
	rows, err := suite.categorySvc.DeleteCategory(context.TODO(), id)
	if err != nil {
		suite.T().Fatal(err)
	}
	zeroRows := int64(0)
	assert.NotEqual(suite.T(), zeroRows, rows)
}
