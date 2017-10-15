package binstruct

import (
	"reflect"
	"testing"

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
