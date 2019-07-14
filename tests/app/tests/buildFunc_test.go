package main

import (
	"testing"

	"github.com/sarulabs/dingo/v3/tests/app/pkg"
	"github.com/sarulabs/dingo/v3/tests/app/generated_services/dic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFunc(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	expected3 := &pkg.BuildFuncTestC{P1: "C"}
	expected2 := pkg.BuildFuncTestB{P1: "B", P2: expected3}
	expected1 := &pkg.BuildFuncTestA{P1: "A", P2: expected2, P3: expected3}
	expected4 := &pkg.BuildFuncTestA{P1: "9999", P2: pkg.BuildFuncTestB{P1: "value", P2: expected3}, P3: expected3}

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
}
