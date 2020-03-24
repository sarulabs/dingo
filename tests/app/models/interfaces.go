package models

// InterfacesTestInterface is an interface used in the tests.
type InterfacesTestInterface interface {
	Test()
}

// InterfacesTestA is a structure used in the tests.
type InterfacesTestA struct {
	Value string
}

// Test allows to implement InterfacesTestInterface.
func (ita InterfacesTestA) Test() {}

// InterfacesTestB is a structure used in the tests.
type InterfacesTestB struct {
	InterfacesTestInterface InterfacesTestInterface
}
