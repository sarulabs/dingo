package services

import (
	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// RetrievalDecls is used in the tests.
var RetrievalDecls = []dingo.Def{
	{
		Name:  "test_retrieval_1",
		Scope: di.App,
		Build: func() (*models.RetrievalTest, error) {
			return models.NewRetrievalTest(), nil
		},
	},
	{
		Name:  "test_retrieval_2",
		Scope: di.Request,
		Build: func() (*models.RetrievalTest, error) {
			return models.NewRetrievalTest(), nil
		},
	},
}
