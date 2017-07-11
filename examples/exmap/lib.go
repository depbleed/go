package exstruct

import "a"

// MyType is an exposed type.
type MyType struct {
	// A is a public map with the key type being provided by a dependency. This is
	// dependency bleeding.
	A map[a.Type]int
	//B is a public map with the value type being provided by a dependency. This is
	// dependency bleeding.
	B map[int]a.Type
	// C is a public map with standard key and value types. Nothing to see there.
	C map[string]int
}
