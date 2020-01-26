package services

import (
	"strconv"

	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// BuildFuncDecls is used in the tests.
var BuildFuncDecls = []dingo.Def{
	{
		Name: "test_build_func_1",
		Build: func(b models.BuildFuncTestB, c *models.BuildFuncTestC) (*models.BuildFuncTestA, error) {
			return &models.BuildFuncTestA{
				P1: "A",
				P2: b,
				P3: c,
			}, nil
		},
	},
	{
		Name: "test_build_func_2",
		Build: func(c *models.BuildFuncTestC) (models.BuildFuncTestB, error) {
			return models.BuildFuncTestB{
				P1: "B",
				P2: c,
			}, nil
		},
	},
	{
		Name: "test_build_func_3",
		Build: func() (*models.BuildFuncTestC, error) {
			return &models.BuildFuncTestC{
				P1: "C",
			}, nil
		},
	},
	{
		Name: "test_build_func_4",
		Build: func(i int, c *models.BuildFuncTestC, s string) (*models.BuildFuncTestA, error) {
			return &models.BuildFuncTestA{
				P1: strconv.Itoa(i),
				P2: models.BuildFuncTestB{
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
