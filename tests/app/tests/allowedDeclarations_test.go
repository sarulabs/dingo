package main

import (
	"testing"

	"github.com/sarulabs/dingo/v4/tests/app/generated/dic"
	"github.com/sarulabs/dingo/v4/tests/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllowedDeclarations(t *testing.T) {
	var obj *models.DeclTypeTest
	var err error

	container, err := dic.NewContainer()
	require.Nil(t, err)

	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType0())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType1())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType2())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType3())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType4())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType5())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType6())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType7())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType8())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType9())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType10())
	assert.Equal(t, &models.DeclTypeTest{}, container.GetTestDeclType11())

	obj, err = container.SafeGetTestDeclType0()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType1()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType2()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType3()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType4()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType5()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType6()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType7()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType8()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType9()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType10()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)

	obj, err = container.SafeGetTestDeclType11()
	assert.Nil(t, err)
	assert.Equal(t, &models.DeclTypeTest{}, obj)
}
