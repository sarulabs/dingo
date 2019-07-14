package main

import (
	"testing"

	"github.com/sarulabs/dingo/v3/tests/app/pkg"
	"github.com/sarulabs/dingo/v3/tests/app/generated_services/dic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClose(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	assert.Equal(t, &pkg.CloseTest{}, container.GetTestClose1())

	res, err := container.SafeGetTestClose1()
	assert.Nil(t, err)
	assert.Equal(t, &pkg.CloseTest{}, res)

	container.Delete()

	assert.Equal(t, &pkg.CloseTest{Closed: true}, res)

	_, err = container.SafeGetTestClose1()
	assert.NotNil(t, err)
}
