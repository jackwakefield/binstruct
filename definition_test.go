package binstruct

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestParseStruct(t *testing.T) {
	foo := struct {
		A int32 `binstruct:"align"`
		B int64
	}{}
	definition, err := parseStruct(foo)
	assert.NotNil(t, definition)
	assert.Equal(t, resolveType(foo), definition.Type)
	assert.Equal(t, 2, len(definition.Fields))
	assert.True(t, definition.HasField("A"))
	assert.True(t, definition.HasField("B"))
	assert.False(t, definition.HasField("C"))
	assert.True(t, definition.HasFieldWithKind("A", reflect.Int32))
	assert.False(t, definition.HasFieldWithKind("A", reflect.Int))
	assert.NoError(t, err)
}

func TestParseStructMissingLenFieldReference(t *testing.T) {
	foo := struct {
		A int32
		B string `binstruct:"lenfield=C"`
	}{}
	_, err := parseStruct(foo)
	assert.EqualError(t, err, "cannot use field C for lenfield")
}

func TestParseStructInvalidLenFieldReference(t *testing.T) {
	foo := struct {
		A string
		B string `binstruct:"lenfield=A"`
	}{}
	_, err := parseStruct(foo)
	assert.EqualError(t, err, "cannot use field A for lenfield")
}

func TestParseStructMissingOffsetReference(t *testing.T) {
	foo := struct {
		A int32
		B string `binstruct:"offsetfield=C"`
	}{}
	_, err := parseStruct(foo)
	assert.EqualError(t, err, "cannot use field C for offsetfield")
}

func TestParseStructInvalidOffsetFieldReference(t *testing.T) {
	foo := struct {
		A string
		B string `binstruct:"offsetfield=A"`
	}{}
	_, err := parseStruct(foo)
	assert.EqualError(t, err, "cannot use field A for offsetfield")
}

func TestParseStructInvalidTagValue(t *testing.T) {
	foo := struct {
		A string `binstruct:"offset=X"`
	}{}
	_, err := parseStruct(foo)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagInt64.Error())
}
