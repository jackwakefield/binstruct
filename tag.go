package binstruct

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type tag map[string]interface{}

var (
	ErrInvalidTagInt64   = errors.New("expected tag value to be a parsable int64")
	ErrInvalidTagFloat64 = errors.New("expected tag value to be a parsable float64")
	ErrInvalidTagString  = errors.New("expected tag value to be a parsable string")
	ErrInvalidTagBool    = errors.New("expected tag value to be a parsable boolean")
)

func parseTag(t reflect.StructTag) tag {
	result := make(map[string]interface{})
	value, _ := t.Lookup("binstruct")
	if value != "" {
		keys := strings.Split(value, ",")
		for _, key := range keys {
			if key != "" {
				if separator := strings.Index(key, "="); separator >= 0 {
					value := parseTagValue(key[separator+1:])
					key = key[0:separator]
					result[key] = value
				} else {
					result[key] = true
				}
			}
		}
	}
	return result
}

func parseTagValue(literal string) interface{} {
	if literal == "true" {
		return true
	}
	if literal == "false" {
		return false
	}
	if index := strings.Index(literal, "."); index >= 0 && index < len(literal)-1 {
		if value, err := strconv.ParseFloat(literal, 64); err == nil {
			return value
		}
	}
	if value, err := strconv.ParseInt(literal, 0, 64); err == nil {
		return value
	}
	return literal
}

func (t tag) Int64(key string) (int64, error) {
	value, ok := t[key]
	if ok {
		option, ok := value.(int64)
		if !ok {
			return 0, ErrInvalidTagInt64
		}
		return option, nil
	}
	return 0, nil
}

func (t tag) Float64(key string) (float64, error) {
	value, ok := t[key]
	if ok {
		option, ok := value.(float64)
		if !ok {
			return 0, ErrInvalidTagFloat64
		}
		return option, nil
	}
	return 0, nil
}

func (t tag) String(key string) (string, error) {
	value, ok := t[key]
	if ok {
		option, ok := value.(string)
		if !ok {
			return "", ErrInvalidTagString
		}
		return option, nil
	}
	return "", nil
}

func (t tag) Byte(key string) (byte, error) {
	value, err := t.String(key)
	if err != nil {
		return 0, err
	}
	if len(value) > 0 {
		return value[0], nil
	}
	return 0, nil
}

func (t tag) Bool(key string) (bool, error) {
	value, ok := t[key]
	if ok {
		option, ok := value.(bool)
		if !ok {
			return false, ErrInvalidTagBool
		}
		return option, nil
	}
	return false, nil
}

func (t tag) Contains(key string) bool {
	_, ok := t[key]
	return ok
}
