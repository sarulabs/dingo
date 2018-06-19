package main

import (
	anotherpkgAliasTest "anotherpkg"
	"otherpkg"
	"testing"

	"github.com/sarulabs/dingo/dingo/tests/app/pkg"

	"github.com/sarulabs/dingo/dingo/tests/app/generated_services/dic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVendor(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	expected2 := &otherpkg.StructOtherPkg{Value: "OK"}
	expected1 := &pkg.VendorTest{Value: expected2}
	expected3 := &anotherpkgAliasTest.StructAnotherPkg{Value: "OK"}

	assert.Equal(t, expected1, container.GetTestVendor1())
	assert.Equal(t, expected2, container.GetTestVendor2())
	assert.Equal(t, expected3, container.GetTestVendor3())

	res1, err := container.SafeGetTestVendor1()
	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)

	res2, err := container.SafeGetTestVendor2()
	assert.Nil(t, err)
	assert.Equal(t, expected2, res2)

	res3, err := container.SafeGetTestVendor3()
	assert.Nil(t, err)
	assert.Equal(t, expected3, res3)
}
