package models

import "github.com/sarulabs/dingo/v4/tests/app/models/testinterfaces"

// InterfacesTestA is a structure used in the tests.
type InterfacesTestA struct {
	Value string
}

// Test allows to implement InterfacesTestInterface.
func (ita InterfacesTestA) Test() {}

// InterfacesTestB is a structure used in the tests.
type InterfacesTestB struct {
	InterfacesTestInterface testinterfaces.InterfacesTestInterface
}
