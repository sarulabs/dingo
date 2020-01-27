package main

import (
	"testing"

	"github.com/sarulabs/dingo/v4/tests/app/generated/dic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScope(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	// create sub-containers
	req1, err := container.SubContainer()
	require.Nil(t, err)

	req2, err := container.SubContainer()
	require.Nil(t, err)

	// retrieve objects
	r1o1, err := req1.SafeGetTestScope1()
	require.Nil(t, err)

	r1o2, err := req1.SafeGetTestScope2()
	require.Nil(t, err)

	r2o1, err := req2.SafeGetTestScope1()
	require.Nil(t, err)

	r2o2, err := req2.SafeGetTestScope2()
	require.Nil(t, err)

	// check values
	assert.Equal(t, r1o1, r2o1)
	assert.NotEqual(t, r1o2, r2o2)
	assert.NotEqual(t, r1o1, r1o2)
	assert.NotEqual(t, r2o1, r2o2)
	assert.Equal(t, r1o1, req1.GetTestScope1())
	assert.Equal(t, r1o2, req1.GetTestScope2())
	assert.Equal(t, r2o1, req2.GetTestScope1())
	assert.Equal(t, r2o2, req2.GetTestScope2())
}
