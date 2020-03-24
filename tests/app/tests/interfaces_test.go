package main

import (
	"testing"

	"github.com/sarulabs/dingo/v4/tests/app/generated/dic"
	"github.com/sarulabs/dingo/v4/tests/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterfaces(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	obj1 := container.Get("test_interfaces_1")
	obj1A := obj1.(*models.InterfacesTestB).InterfacesTestInterface.(models.InterfacesTestA)
	assert.Equal(t, "1", obj1A.Value)

	obj2 := container.Get("test_interfaces_2")
	obj2A := obj2.(*models.InterfacesTestB).InterfacesTestInterface.(models.InterfacesTestA)
	assert.Equal(t, "2", obj2A.Value)
}
