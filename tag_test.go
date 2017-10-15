package binstruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	var emptyTag reflect.StructTag = `binstruct:""`
	assert.Equal(t, 0, len(parseTag(emptyTag)))

	var nilTag reflect.StructTag
	assert.Equal(t, 0, len(parseTag(nilTag)))

	var fullTag reflect.StructTag = `binstruct:"a,b=1,c=1.0,d=test,e=false,f=true"`
	parsedTag := parseTag(fullTag)
	assert.Equal(t, 6, len(parsedTag))

	expectedTag := make(tag, 0)
	expectedTag["a"] = true
	expectedTag["b"] = int64(1)
	expectedTag["c"] = 1.0
	expectedTag["d"] = "test"
	expectedTag["e"] = false
	expectedTag["f"] = true
	assert.Equal(t, expectedTag, parsedTag)
}

func TestParseTagValue(t *testing.T) {
	assert.Equal(t, true, parseTagValue("true"))
	assert.Equal(t, false, parseTagValue("false"))
	assert.Equal(t, 0.01, parseTagValue("0.01"))
	assert.Equal(t, "0.", parseTagValue("0."))
	assert.Equal(t, int64(1), parseTagValue("1"))
	assert.Equal(t, "foo", parseTagValue("foo"))
}

func TestParseTagBool(t *testing.T) {
	values := make(tag, 0)
	values["a"] = true

	assert.Equal(t, true, values.Contains("a"))
	a, err := values.Bool("a")
	assert.Equal(t, true, a)
	assert.NoError(t, err)
	_, err = values.Byte("a")
	assert.Error(t, err)
	_, err = values.Int64("a")
	assert.Error(t, err)
	_, err = values.Float64("a")
	assert.Error(t, err)
	_, err = values.String("a")
	assert.Error(t, err)
}

func TestParseTagByte(t *testing.T) {
	values := make(tag, 0)
	values["a"] = "x"

	assert.Equal(t, true, values.Contains("a"))
	_, err := values.Bool("a")
	assert.Error(t, err)
	a, err := values.Byte("a")
	assert.Equal(t, byte('x'), a)
	assert.NoError(t, err)
	_, err = values.Int64("a")
	assert.Error(t, err)
	_, err = values.Float64("a")
	assert.Error(t, err)
	stringValue, err := values.String("a")
	assert.Equal(t, "x", stringValue)
	assert.NoError(t, err)
}

func TestParseTagInt64(t *testing.T) {
	values := make(tag, 0)
	values["a"] = int64(1)

	assert.Equal(t, true, values.Contains("a"))
	_, err := values.Bool("a")
	assert.Error(t, err)
	_, err = values.Byte("a")
	assert.Error(t, err)
	a, err := values.Int64("a")
	assert.Equal(t, int64(1), a)
	assert.NoError(t, err)
	_, err = values.Float64("a")
	assert.Error(t, err)
	_, err = values.String("a")
	assert.Error(t, err)
}

func TestParseTagFloat64(t *testing.T) {
	values := make(tag, 0)
	values["a"] = float64(1.0)

	assert.Equal(t, true, values.Contains("a"))
	_, err := values.Bool("a")
	assert.Error(t, err)
	_, err = values.Byte("a")
	assert.Error(t, err)
	_, err = values.Int64("a")
	assert.Error(t, err)
	a, err := values.Float64("a")
	assert.Equal(t, float64(1.0), a)
	assert.NoError(t, err)
	_, err = values.String("a")
	assert.Error(t, err)
}

func TestParseTagString(t *testing.T) {
	values := make(tag, 0)
	values["a"] = "foo"

	assert.Equal(t, true, values.Contains("a"))
	_, err := values.Bool("a")
	assert.Error(t, err)
	byteValue, err := values.Byte("a")
	assert.Equal(t, byte('f'), byteValue)
	assert.NoError(t, err)
	_, err = values.Int64("a")
	assert.Error(t, err)
	_, err = values.Float64("a")
	assert.Error(t, err)
	a, err := values.String("a")
	assert.Equal(t, "foo", a)
	assert.NoError(t, err)
}

func TestParseTagMissing(t *testing.T) {
	values := make(tag, 0)

	a, err := values.Bool("a")
	assert.Equal(t, false, a)
	assert.NoError(t, err)

	b, err := values.Byte("b")
	assert.Equal(t, byte(0), b)
	assert.NoError(t, err)

	c, err := values.Int64("c")
	assert.Equal(t, int64(0), c)
	assert.NoError(t, err)

	d, err := values.Float64("d")
	assert.Equal(t, float64(0.0), d)
	assert.NoError(t, err)

	e, err := values.String("e")
	assert.Equal(t, "", e)
	assert.NoError(t, err)
}
