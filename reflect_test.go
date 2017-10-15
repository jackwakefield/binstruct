package binstruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltInType(t *testing.T) {
	var a int
	aType := resolveType(a)
	assert.Equal(t, reflect.Int, aType.Kind())

	var b *int
	bType := resolveType(b)
	assert.Equal(t, reflect.Int, bType.Kind())
}

func TestStructType(t *testing.T) {
	type foo struct {
	}

	var a foo
	aType := resolveType(a)
	assert.Equal(t, reflect.Struct, aType.Kind())
	assert.Equal(t, "foo", aType.Name())

	var b *foo
	bType := resolveType(b)
	assert.Equal(t, reflect.Struct, bType.Kind())
	assert.Equal(t, "foo", bType.Name())

	c := (interface{})(a)
	cType := resolveType(c)
	assert.Equal(t, reflect.Struct, cType.Kind())
	assert.Equal(t, "foo", cType.Name())

	d := (interface{})(b)
	dType := resolveType(d)
	assert.Equal(t, reflect.Struct, dType.Kind())
	assert.Equal(t, "foo", dType.Name())
}
