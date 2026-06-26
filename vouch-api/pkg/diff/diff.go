// Package diff provides helpers for detecting changes between two structs.
package diff

import "reflect"

// Field represents a changed field between two values.
type Field struct {
	Name string
	Old  any
	New  any
}

// Struct compares two structs of the same type and returns a list of changed fields.
// Only exported fields are compared. Panics if a and b are not structs.
func Struct(a, b any) []Field {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}
	if vb.Kind() == reflect.Ptr {
		vb = vb.Elem()
	}

	t := va.Type()
	var changes []Field

	for i := 0; i < va.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fa := va.Field(i).Interface()
		fb := vb.Field(i).Interface()
		if !reflect.DeepEqual(fa, fb) {
			changes = append(changes, Field{Name: field.Name, Old: fa, New: fb})
		}
	}
	return changes
}

// HasChanges returns true if any exported fields differ between a and b.
func HasChanges(a, b any) bool {
	return len(Struct(a, b)) > 0
}

// FieldNames returns just the names of changed fields.
func FieldNames(a, b any) []string {
	fields := Struct(a, b)
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names
}
