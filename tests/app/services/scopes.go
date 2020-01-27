package services

import (
	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// ScopeDecls is used in the tests.
var ScopeDecls = []dingo.Def{
	{
		Name:  "test_scope_1",
		Scope: di.App,
		Build: func() (*models.ScopeTest, error) {
			return models.NewScopeTest(), nil
		},
	},
	{
		Name:  "test_scope_2",
		Scope: di.Request,
		Build: func() (*models.ScopeTest, error) {
			return models.NewScopeTest(), nil
		},
	},
}
