package main

import (
	"testing"

	"github.com/sarulabs/dingo/v3/tests/app/pkg"
	"github.com/sarulabs/dingo/v3/tests/app/generated_services/dic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAutofill(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	expected1 := &pkg.AutofillTestA{Value: "A1"}
	expected2 := &pkg.AutofillTestA{Value: "A2"}
	expected3 := &pkg.AutofillTestB{Value: expected2}

	assert.Equal(t, expected1, container.GetTestAutofill1())
	assert.Equal(t, expected2, container.GetTestAutofill2())
	assert.Equal(t, expected3, container.GetTestAutofill3())

	res1, err := container.SafeGetTestAutofill1()
	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)

	res2, err := container.SafeGetTestAutofill2()
	assert.Nil(t, err)
	assert.Equal(t, expected2, res2)

	res3, err := container.SafeGetTestAutofill3()
	assert.Nil(t, err)
	assert.Equal(t, expected3, res3)
}
