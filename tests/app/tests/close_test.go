package main

import (
	"testing"

	"github.com/sarulabs/dingo/v4/tests/app/generated/dic"
	"github.com/sarulabs/dingo/v4/tests/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClose(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	assert.Equal(t, &models.CloseTest{}, container.GetTestClose1())

	res, err := container.SafeGetTestClose1()
	assert.Nil(t, err)
	assert.Equal(t, &models.CloseTest{}, res)

	container.Delete()

	assert.Equal(t, &models.CloseTest{Closed: true}, res)

	_, err = container.SafeGetTestClose1()
	assert.NotNil(t, err)
}
