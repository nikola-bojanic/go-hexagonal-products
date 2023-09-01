package category

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
	categoryHttpSvc CategoryHttpHandler
	wsContainer     *restful.Container
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
	suite.categoryHttpSvc = *NewCategoryHandler(realCategorySvc, suite.wsContainer)

}

func TestCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(HttpSuite))
}
func (suite *HttpSuite) TestGetCategories() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categoryHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "GET", "/category", nil, nil)
	var response []CategoryModel
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling category response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.NotNil(suite.T(), response)
	suite.categoryHttpSvc.categorySvc.DeleteCategory(context.TODO(), id)
}
func (suite *HttpSuite) TestGetCategory() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categoryHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "GET", "/category/"+strconv.Itoa(int(id)), nil, nil)
	var response CategoryModel
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling category response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.Equal(suite.T(), testCategory.Name, response.Name)
	suite.categoryHttpSvc.categorySvc.DeleteCategory(context.TODO(), id)

}
func (suite *HttpSuite) TestCreateCategory() {
	testCategory := CategoryRequest{
		Name: "test",
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "POST", "/category", testCategory, nil)
	var createResponse Response
	err := json.Unmarshal(responseRec.Body.Bytes(), &createResponse)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling category response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	assert.Equal(suite.T(), testCategory.Name, createResponse.Name)

	suite.categoryHttpSvc.categorySvc.DeleteCategory(context.TODO(), createResponse.ID)
}
func (suite *HttpSuite) TestUpdateCategory() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categoryHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	updateCategory := CategoryRequest{
		Name: "testing",
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "PUT", "/category/"+strconv.Itoa(int(id)), updateCategory, nil)
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	var updateResponse Response
	err = json.Unmarshal(responseRec.Body.Bytes(), &updateResponse)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling category response: %s", err)
	}
	assert.Equal(suite.T(), updateCategory.Name, updateResponse.Name)
	rowsAffected := int64(1)
	assert.Equal(suite.T(), rowsAffected, updateResponse.ID)
	suite.categoryHttpSvc.categorySvc.DeleteCategory(context.TODO(), id)
}
func (suite *HttpSuite) TestDeleteCategory() {
	testCategory := domain.Category{
		Name: "test",
	}
	id, err := suite.categoryHttpSvc.categorySvc.CreateCategory(context.TODO(), &testCategory)
	if err != nil {
		suite.T().Fatalf("Error creating test category: %s", err)
	}
	responseRec := testutil.MakeRequest(suite.wsContainer, "DELETE", "/category/"+strconv.Itoa(int(id)), nil, nil)
	var response Response
	err = json.Unmarshal(responseRec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatalf("Error unmarshalling category response: %s", err)
	}
	assert.Equal(suite.T(), http.StatusOK, responseRec.Code)
	message := "category deleted"
	rows := int64(1)
	assert.Equal(suite.T(), message, response.Name)
	assert.Equal(suite.T(), rows, response.ID)
}
