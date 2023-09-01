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

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type UserSuite struct {
	suite.Suite
	userRep *repo.UserRepository
	userSvc *UserService
}

func (suite *UserSuite) SetupTest() {
}

func (suite *UserSuite) TearDownTest() {
	// todo: clear the DB
}

func (suite *UserSuite) SetupSuite() {

	app := testutil.InitTestApp()
	suite.userRep = repo.NewUserRepository(app.DB)
	suite.userSvc = NewUserService(suite.userRep)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) TestUserRegistration() {

	userEmail := "email1@provider.com"
	err := suite.userSvc.RegisterUser(context.TODO(), &domain.User{Email: userEmail})

	if err != nil {
		suite.T().Fatal(err)
	}

	user, err := suite.userSvc.FindByEmail(context.TODO(), userEmail)
	if err != nil {
		suite.T().Fatal(err)
	}

	assert.Equal(suite.T(), user.Email, userEmail)
}

func (suite *UserSuite) TestCannotRegisterWithExistingEmail() {
	userEmail := "email2@provider.com"
	err := suite.userSvc.RegisterUser(context.TODO(), &domain.User{Email: userEmail})

	if err != nil {
		suite.T().Fatal(err)
	}

	err = suite.userSvc.RegisterUser(context.TODO(), &domain.User{Email: userEmail})

	assert.ErrorIs(suite.T(), err, repo.ErrDuplicateEmail)
}
