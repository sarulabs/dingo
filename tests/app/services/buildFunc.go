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
	{
		Name: "test_build_func_5",
		Build: func() (models.TypeBasedOnBasicType, error) {
			return models.TypeBasedOnBasicType(999), nil
		},
	},
	{
		Name: "test_build_func_6",
		Build: func() (models.TypeBasedOnSliceOfBasicType, error) {
			return models.TypeBasedOnSliceOfBasicType([]byte("test")), nil
		},
	},
	{
		Name: "test_build_func_7",
		Build: func() (struct{}, error) {
			return struct{}{}, nil
		},
	},
	{
		Name: "test_build_func_8",
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
		Params: dingo.NewFuncParams(
			9999,
			&models.BuildFuncTestC{P1: "C"},
			"value",
		),
	},
}
