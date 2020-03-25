package services

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
	"github.com/sarulabs/dingo/v4/tests/app/models/testinterfaces"
)

// InterfacesDecls is used in the tests.
var InterfacesDecls = []dingo.Def{
	{
		Name: "test_interfaces_1",
		Build: func(i testinterfaces.InterfacesTestInterface) (*models.InterfacesTestB, error) {
			return &models.InterfacesTestB{InterfacesTestInterface: i}, nil
		},
		Params: dingo.Params{
			"0": models.InterfacesTestA{Value: "1"},
		},
	},
	{
		Name:  "test_interfaces_2",
		Build: (*models.InterfacesTestB)(nil),
		Params: dingo.Params{
			"InterfacesTestInterface": models.InterfacesTestA{Value: "2"},
		},
	},
}
