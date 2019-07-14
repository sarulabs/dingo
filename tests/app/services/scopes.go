package services

import (
	"github.com/sarulabs/dingo/v2"
	"github.com/sarulabs/dingo/v2/tests/app/pkg"
)

// ScopeDecls is used in the tests.
var ScopeDecls = []dingo.Def{
	{
		Name:  "test_scope_1",
		Scope: dingo.App,
		Build: func() (*pkg.ScopeTest, error) {
			return pkg.NewScopeTest(), nil
		},
	},
	{
		Name:  "test_scope_2",
		Scope: dingo.Request,
		Build: func() (*pkg.ScopeTest, error) {
			return pkg.NewScopeTest(), nil
		},
	},
}
