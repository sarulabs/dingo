package services

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// UnsharedDecls is used in the tests.
var UnsharedDecls = []dingo.Def{
	{
		Name: "test_unshared_1",
		Build: func() (*models.UnsharedTest, error) {
			return models.NewUnsharedTest(), nil
		},
		Unshared: true,
	},
	{
		Name: "test_unshared_2",
		Build: func() (*models.UnsharedTest, error) {
			return models.NewUnsharedTest(), nil
		},
		Unshared: false,
	},
}
