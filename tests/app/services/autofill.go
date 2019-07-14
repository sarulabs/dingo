package services

import (
	"github.com/sarulabs/dingo/v2"
	"github.com/sarulabs/dingo/v2/tests/app/pkg"
)

// AutofillDecls is used in the tests.
var AutofillDecls = []dingo.Def{
	{
		Name:           "test_autofill_1",
		Build:          (*pkg.AutofillTestA)(nil),
		NotForAutoFill: true,
		Params:         dingo.Params{"Value": "A1"},
	},
	{
		Name:   "test_autofill_2",
		Build:  (*pkg.AutofillTestA)(nil),
		Params: dingo.Params{"Value": "A2"},
	},
	{
		Name:  "test_autofill_3",
		Build: (*pkg.AutofillTestB)(nil),
	},
}
