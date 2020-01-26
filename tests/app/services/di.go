package services

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// DiDecls is used in the tests.
var DiDecls = []dingo.Def{
	{
		Name: "test_di_1",
		Build: func() (models.DiTest, error) {
			return models.DiTest{Value: "1"}, nil
		},
	},
	{
		Name: "test_di_2",
		Build: func() (models.DiTest, error) {
			return models.DiTest{Value: "2"}, nil
		},
	},
	{
		Name: "test_di_3",
		Build: func() (models.DiTest, error) {
			return models.DiTest{Value: "3"}, nil
		},
	},
}
