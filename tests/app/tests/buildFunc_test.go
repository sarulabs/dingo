package main

import (
	"testing"

	"github.com/sarulabs/dingo/v4/tests/app/generated/dic"
	"github.com/sarulabs/dingo/v4/tests/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFunc(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	expected3 := &models.BuildFuncTestC{P1: "C"}
	expected2 := models.BuildFuncTestB{P1: "B", P2: expected3}
	expected1 := &models.BuildFuncTestA{P1: "A", P2: expected2, P3: expected3}
	expected4 := &models.BuildFuncTestA{P1: "9999", P2: models.BuildFuncTestB{P1: "value", P2: expected3}, P3: expected3}
	expected5 := models.TypeBasedOnBasicType(999)
	expected6 := models.TypeBasedOnSliceOfBasicType([]byte("test"))

	assert.Equal(t, expected1, container.GetTestBuildFunc1())
	assert.Equal(t, expected2, container.GetTestBuildFunc2())
	assert.Equal(t, expected3, container.GetTestBuildFunc3())
	assert.Equal(t, expected4, container.GetTestBuildFunc4())

	res1, err := container.SafeGetTestBuildFunc1()
	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)

	res2, err := container.SafeGetTestBuildFunc2()
	assert.Nil(t, err)
	assert.Equal(t, expected2, res2)

	res3, err := container.SafeGetTestBuildFunc3()
	assert.Nil(t, err)
	assert.Equal(t, expected3, res3)

	res4, err := container.SafeGetTestBuildFunc4()
	assert.Nil(t, err)
	assert.Equal(t, expected4, res4)

	res5, err := container.SafeGetTestBuildFunc5()
	assert.Nil(t, err)
	assert.Equal(t, expected5, res5)

	res6, err := container.SafeGetTestBuildFunc6()
	assert.Nil(t, err)
	assert.Equal(t, expected6, res6)

	res7, err := container.SafeGetTestBuildFunc7()
	assert.Nil(t, err)
	assert.Equal(t, struct{}{}, res7)
}
