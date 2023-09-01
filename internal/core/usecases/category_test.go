package usecases

import (
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/stretchr/testify/suite"
)

type CategorySuite struct {
	suite.Suite
	categoryRep *repo.CategoryRepository
	categorySvc *CategoryService
}
