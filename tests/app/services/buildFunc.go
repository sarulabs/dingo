package services

import (
	"strconv"

	"github.com/sarulabs/dingo/v2"
	"github.com/sarulabs/dingo/v2/tests/app/pkg"
)

// BuildFuncDecls is used in the tests.
var BuildFuncDecls = []dingo.Def{
	{
		Name: "test_build_func_1",
		Build: func(b pkg.BuildFuncTestB, c *pkg.BuildFuncTestC) (*pkg.BuildFuncTestA, error) {
			return &pkg.BuildFuncTestA{
				P1: "A",
				P2: b,
				P3: c,
			}, nil
		},
	},
	{
		Name: "test_build_func_2",
		Build: func(c *pkg.BuildFuncTestC) (pkg.BuildFuncTestB, error) {
			return pkg.BuildFuncTestB{
				P1: "B",
				P2: c,
			}, nil
		},
	},
	{
		Name: "test_build_func_3",
		Build: func() (*pkg.BuildFuncTestC, error) {
			return &pkg.BuildFuncTestC{
				P1: "C",
			}, nil
		},
	},
	{
		Name: "test_build_func_4",
		Build: func(i int, c *pkg.BuildFuncTestC, s string) (*pkg.BuildFuncTestA, error) {
			return &pkg.BuildFuncTestA{
				P1: strconv.Itoa(i),
				P2: pkg.BuildFuncTestB{
					P1: s,
					P2: c,
				},
				P3: c,
			}, nil
		},
		Params: dingo.Params{
			"0": 9999,
			"2": "value",
		},
	},
}
