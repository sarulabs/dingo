package services

import (
	"github.com/sarulabs/dingo/v2"
	"github.com/sarulabs/dingo/v2/tests/app/pkg"
)

// CloseDecls is used in the tests.
var CloseDecls = []dingo.Def{
	{
		Name:  "test_close_1",
		Build: (*pkg.CloseTest)(nil),
		Close: func(ct *pkg.CloseTest) error {
			ct.Closed = true
			return nil
		},
	},
}
