package binstruct

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaultOptions(t *testing.T) {
	options := &FieldOptions{}
	*options = *defaultFieldOptions
	options.Align = true
	SetDefaultOptions(options)
	assert.Equal(t, defaultFieldOptions, options)
}

func TestParseTagFieldOptions(t *testing.T) {
	options, err := parseTagFieldOptions(nil)
	assert.Equal(t, defaultFieldOptions, options)
	assert.NoError(t, err)

	var fullTag reflect.StructTag = `binstruct:"skip=-1,offset=1,offsetfield=foo,len=2,lenfield=bar,stringtype=null,stringpad=b,align,alignbytes=8,mask=0xFFFFFFFF"`
	tag := parseTag(fullTag)
	options, err = parseTagFieldOptions(tag)
	assert.NoError(t, err)
	assert.Equal(t, &FieldOptions{
		Skip:        -1,
		Offset:      int64(1),
		OffsetField: "foo",
		Len:         int64(2),
		LenField:    "bar",
		StringType:  StringNullTerminated,
		StringPad:   byte('b'),
		Align:       true,
		AlignBytes:  8,
		Mask:        0xFFFFFFFF,
	}, options)
}

func TestParseTagFieldInvalidSkip(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"skip=A"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagInt64.Error())
}

func TestParseTagFieldInvalidOffset(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"offset=A"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagInt64.Error())
}

func TestParseTagFieldInvalidOffsetField(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"offsetfield"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagString.Error())
}

func TestParseTagFieldInvalidLen(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"len=A"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagInt64.Error())
}

func TestParseTagFieldInvalidLenField(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"lenfield"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagString.Error())
}

func TestParseTagFieldInvalidStringType(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"stringtype"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagString.Error())
}

func TestParseTagFieldInvalidStringPad(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"stringpad"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagString.Error())
}

func TestParseTagFieldInvalidAlign(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"align=1"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagBool.Error())
}

func TestParseTagFieldInvalidAlignBytes(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"alignbytes=false"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagInt64.Error())
}

func TestParseTagFieldInvalidMask(t *testing.T) {
	tag := parseTag(reflect.StructTag(`binstruct:"mask=false"`))
	_, err := parseTagFieldOptions(tag)
	assert.EqualError(t, errors.Cause(err), ErrInvalidTagInt64.Error())
}
