package models

import "fmt"

// Type represents a field that isn't contained.
type Type struct {
	Field *Field // The field information for the type.
}

//nolint:unused // isStruct returns whether the type is a struct.
func (t Type) isStruct() bool {
	return t.Field.Definition == "struct"
}

//nolint:unused // isInterface returns whether the type is an interface.
func (t Type) isInterface() bool {
	return t.Field.Definition == "interface"
}

// ParameterName gets the parameter name of the type.
func (t Type) ParameterName() string {
	return t.Field.Pointer + t.Field.Definition
}

func (t Type) String() string {
	return fmt.Sprintf("type %v", t.Field.FullName(""))
}
