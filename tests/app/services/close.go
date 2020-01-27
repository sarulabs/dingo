package services

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// CloseDecls is used in the tests.
var CloseDecls = []dingo.Def{
	{
		Name:  "test_close_1",
		Build: (*models.CloseTest)(nil),
		Close: func(ct *models.CloseTest) error {
			ct.Closed = true
			return nil
		},
	},
}
