package services

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// AutofillDecls is used in the tests.
var AutofillDecls = []dingo.Def{
	{
		Name:           "test_autofill_1",
		Build:          (*models.AutofillTestA)(nil),
		NotForAutoFill: true,
		Params:         dingo.Params{"Value": "A1"},
		Description: `Test description.

Even on multiple lines.`,
	},
	{
		Name:   "test_autofill_2",
		Build:  (*models.AutofillTestA)(nil),
		Params: dingo.Params{"Value": "A2"},
	},
	{
		Name:  "test_autofill_3",
		Build: (*models.AutofillTestB)(nil),
	},
}
