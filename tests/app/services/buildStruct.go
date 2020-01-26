package services

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// BuildStructDecls is used in the tests.
var BuildStructDecls = []dingo.Def{
	{
		Name:  "test_build_struct_1",
		Build: (*models.BuildStructTestA)(nil),
		Params: dingo.Params{
			"P1": "A",
		},
	},
	{
		Name:  "test_build_struct_2",
		Build: (*models.BuildStructTestB)(nil),
		Params: dingo.Params{
			"P1": "B",
		},
	},
	{
		Name:  "test_build_struct_3",
		Build: (*models.BuildStructTestC)(nil),
		Params: dingo.Params{
			"P1": "C",
		},
	},
	{
		Name:  "test_build_struct_4",
		Build: (*models.BuildStructTestA)(nil),
		Params: dingo.Params{
			"P1": "value1",
			"P3": &models.BuildStructTestC{P1: "value2"},
		},
	},
}
