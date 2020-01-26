package main

import (
	"testing"

	"github.com/sarulabs/dingo/v4/tests/app/generated/dic"
	"github.com/sarulabs/dingo/v4/tests/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDi(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	obj1 := container.Get("test_di_1")
	assert.Equal(t, "1", obj1.(models.DiTest).Value)

	obj2, err := container.SafeGet("test_di_2")
	assert.Nil(t, err)
	assert.Equal(t, "2", obj2.(models.DiTest).Value)

	var obj3 models.DiTest
	err = container.Fill("test_di_3", &obj3)
	assert.Nil(t, err)
	assert.Equal(t, "3", obj3.Value)
}
