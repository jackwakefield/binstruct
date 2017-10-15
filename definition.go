package binstruct

import (
	"fmt"
	"reflect"
)

type structDefinition struct {
	Type   reflect.Type
	Fields []*fieldDefinition
}

// parseStruct creates a struct definition from the struct type.
func parseStruct(value interface{}) (*structDefinition, error) {
	t := resolveType(value)
	return parseStructType(t)
}

// parseStructType creates a struct definition from the type detailing the fields and their options.
func parseStructType(t reflect.Type) (*structDefinition, error) {
	definition := &structDefinition{Type: t}
	if err := definition.parseFields(); err != nil {
		return nil, err
	}
	return definition, nil
}

// parseFields recursively iterates through the struct's fields creating fieldDefinition.
func (s *structDefinition) parseFields() error {
	fieldCount := s.Type.NumField()
	fieldDefinitions := make([]*fieldDefinition, fieldCount)
	for i := 0; i < fieldCount; i++ {
		field := s.Type.Field(i)
		var err error
		// attempt to parse the field
		if fieldDefinitions[i], err = parseField(s, field); err != nil {
			return err
		}
	}
	s.Fields = fieldDefinitions
	return nil
}

// HasField determines whether the field name exists.
func (s *structDefinition) HasField(name string) bool {
	for _, field := range s.Fields {
		if field.Field.Name == name {
			return true
		}
	}
	return false
}

// HasFieldWithKind determines whether the field name exists, and whether the kind is equal to any of the values given.
func (s *structDefinition) HasFieldWithKind(name string, kinds ...reflect.Kind) bool {
	for _, field := range s.Fields {
		if field.Field.Name == name {
			// iterate through and check against each kind given
			for _, kind := range kinds {
				if field.Field.Type.Kind() == kind {
					return true
				}
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
	Children []*structDefinition
}

// numericalFieldKinds contains a list of the valid kinds when dealing with numbers (offsets, lengths etc.)
var numericalFieldKinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
}

// parseField creates a field definition from the given field type belonging to the struct.
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
		return err
	}

	// ensure the options referencing other fields are valid kinds
	if f.Options.LenField != "" {
		if !f.Struct.HasFieldWithKind(f.Options.LenField, numericalFieldKinds...) {
			return fmt.Errorf("cannot use field %s for lenfield", f.Options.LenField)
		}
	}
	if f.Options.OffsetField != "" {
		if !f.Struct.HasFieldWithKind(f.Options.OffsetField, numericalFieldKinds...) {
			return fmt.Errorf("cannot use field %s for offsetfield", f.Options.OffsetField)
		}
	}

	return nil
}
