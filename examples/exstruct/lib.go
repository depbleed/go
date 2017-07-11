package exstruct

import "a"

// MyType is an exposed type.
type MyType struct {
	// A is a public member of a type provided by a dependency. This is
	// dependency bleeding.
	A a.Type
}
