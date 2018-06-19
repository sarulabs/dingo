package main

import (
	"testing"

	"github.com/sarulabs/dingo/dingo/tests/app/generated_services/dic"
	"github.com/stretchr/testify/require"
)

func TestRetrieval(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	// unscoped retrieval
	o1, err := container.UnscopedSafeGetTestRetrieval2()
	require.Nil(t, err)

	o2, err := container.UnscopedSafeGetTestRetrieval2()
	require.Nil(t, err)

	require.Equal(t, o1, o2)

	container.Clean()

	o3, err := container.UnscopedSafeGetTestRetrieval2()
	require.Nil(t, err)

	require.NotEqual(t, o1, o3)

	// test retrieval functions
	o4 := dic.TestRetrieval1(container)
	require.NotNil(t, o4)
}
