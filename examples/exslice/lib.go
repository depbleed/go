package exstruct

import "a"

// MyType is an exposed type.
type MyType struct {
	// A is a public member of a slice type provided by a dependency. This is
	// dependency bleeding.
	A []a.Type
	B []int
}
