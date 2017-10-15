package binstruct

import (
	"reflect"

	"github.com/pkg/errors"
)

type structDefinition struct {
	Type   reflect.Type
	Fields map[string]*fieldDefinition
}

// parseStruct creates a struct definition from the struct type.
func parseStruct(value interface{}) (*structDefinition, error) {
	t := resolveType(value)
	return parseStructType(t)
}

// parseStructType creates a struct definition from the type
// detailing the fields and their options.
func parseStructType(t reflect.Type) (*structDefinition, error) {
	definition := &structDefinition{Type: t}
	if err := definition.parseFields(); err != nil {
		return nil, err
	}
	return definition, nil
}

// parseFields recursively iterates through the struct's fields
// creating fieldDefinition.
func (s *structDefinition) parseFields() error {
	fieldCount := s.Type.NumField()
	fieldDefinitions := make(map[string]*fieldDefinition, fieldCount)
	for i := 0; i < fieldCount; i++ {
		field := s.Type.Field(i)

		var err error
		// attempt to parse the field
		if fieldDefinitions[field.Name], err = parseField(s, field); err != nil {
			return err
		}
	}
	s.Fields = fieldDefinitions
	return nil
}

// HasField determines whether the field name exists.
func (s *structDefinition) HasField(name string) bool {
	_, ok := s.Fields[name]
	return ok
}

// HasFieldWithKind determines whether the field name exists and
// the type-kind is equal to any of those given.
func (s *structDefinition) HasFieldWithKind(name string, kinds ...reflect.Kind) bool {
	if field, ok := s.Fields[name]; ok {
		for _, kind := range kinds {
			if field.Field.Type.Kind() == kind {
				return true
			}
		}
	}
	return false
}

type fieldDefinition struct {
	Struct   *structDefinition
	Field    reflect.StructField
	Type     reflect.Type
	Options  *FieldOptions
	Children *structDefinition
}

// numericalFieldKinds contains a list of the valid kinds when
// dealing with numbers (offsets, lengths etc.)
var numericalFieldKinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
}

var (
	ErrTagParseFailed           = errors.New("failed to parse field tag")
	ErrOptionLenFieldInvalid    = errors.New("tag option lenfield must be an integer")
	ErrOptionOffsetFieldInvalid = errors.New("tag option offsetfield must be an integer")
)

// parseField creates a field definition from the given
// field type belonging to the struct.
func parseField(s *structDefinition, field reflect.StructField) (*fieldDefinition, error) {
	definition := &fieldDefinition{
		Struct: s,
		Field:  field,
		Type:   underlyingType(field.Type),
	}
	if err := definition.parseTag(); err != nil {
		return nil, err
	}
	return definition, nil
}

// parseTag retrieves the field options from the struct tag.
func (f *fieldDefinition) parseTag() error {
	tag := parseTag(f.Field.Tag)
	var err error
	if f.Options, err = parseTagFieldOptions(tag); err != nil {
		return errors.Wrap(err, ErrTagParseFailed.Error())
	}

	// ensure the options referencing other fields exist and are valid
	if f.Options.LenField != "" {
		if !f.Struct.HasFieldWithKind(f.Options.LenField, numericalFieldKinds...) {
			return ErrOptionLenFieldInvalid
		}
	}
	if f.Options.OffsetField != "" {
		if !f.Struct.HasFieldWithKind(f.Options.OffsetField, numericalFieldKinds...) {
			return ErrOptionOffsetFieldInvalid
		}
	}

	return nil
}
