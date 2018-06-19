package services

import (
	"github.com/sarulabs/dingo"
	"github.com/sarulabs/dingo/dingo/tests/app/pkg"
)

// RetrievalDecls is used in the tests.
var RetrievalDecls = []dingo.Def{
	{
		Name:  "test_retrieval_1",
		Scope: dingo.App,
		Build: func() (*pkg.RetrievalTest, error) {
			return pkg.NewRetrievalTest(), nil
		},
	},
	{
		Name:  "test_retrieval_2",
		Scope: dingo.Request,
		Build: func() (*pkg.RetrievalTest, error) {
			return pkg.NewRetrievalTest(), nil
		},
	},
}
