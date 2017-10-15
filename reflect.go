package binstruct

import "reflect"

// resolveType resolves the underlying type of the given value.
func resolveType(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)
	return underlyingType(t)
}

// underlyingType resolves the underlying type by iterating through
// the current and parent types until a kind is found which is not a
// pointer of interface.
func underlyingType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	return t
}
