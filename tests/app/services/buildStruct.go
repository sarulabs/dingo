package services

import (
	"github.com/sarulabs/dingo/v2"
	"github.com/sarulabs/dingo/v2/tests/app/pkg"
)

// BuildStructDecls is used in the tests.
var BuildStructDecls = []dingo.Def{
	{
		Name:  "test_build_struct_1",
		Build: (*pkg.BuildStructTestA)(nil),
		Params: dingo.Params{
			"P1": "A",
		},
	},
	{
		Name:  "test_build_struct_2",
		Build: (*pkg.BuildStructTestB)(nil),
		Params: dingo.Params{
			"P1": "B",
		},
	},
	{
		Name:  "test_build_struct_3",
		Build: (*pkg.BuildStructTestC)(nil),
		Params: dingo.Params{
			"P1": "C",
		},
	},
	{
		Name:  "test_build_struct_4",
		Build: (*pkg.BuildStructTestA)(nil),
		Params: dingo.Params{
			"P1": "value1",
			"P3": &pkg.BuildStructTestC{P1: "value2"},
		},
	},
}
