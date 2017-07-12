package expointer

import "a"

// MyType is an exposed type.
type MyType struct {
	// A is a public member of a slice type provided by a dependency. This is
	// dependency bleeding.
	A *a.Type
	// B is a slice of a standard type. Nothing to see here.
	B *int
}
