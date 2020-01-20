package services

import (
	"github.com/sarulabs/dingo/v3"
	"github.com/sarulabs/dingo/v3/tests/app/pkg"
)

// DiDecls is used in the tests.
var DiDecls = []dingo.Def{
	{
		Name: "test_di_1",
		Build: func() (pkg.DiTest, error) {
			return pkg.DiTest{Value: "1"}, nil
		},
	},
	{
		Name: "test_di_2",
		Build: func() (pkg.DiTest, error) {
			return pkg.DiTest{Value: "2"}, nil
		},
	},
	{
		Name: "test_di_3",
		Build: func() (pkg.DiTest, error) {
			return pkg.DiTest{Value: "3"}, nil
		},
	},
}
