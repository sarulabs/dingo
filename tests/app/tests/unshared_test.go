package main

import (
	"testing"

	"github.com/sarulabs/dingo/v4/tests/app/generated/dic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnshared(t *testing.T) {
	container, err := dic.NewContainer()
	require.Nil(t, err)

	// Unshared true
	assert.NotEqual(t, container.GetTestUnshared1(), container.GetTestUnshared1())

	// Unshared false
	assert.Equal(t, container.GetTestUnshared2(), container.GetTestUnshared2())
}
