package main

import "a"

// MyType is an exposed type.
type MyType struct {
	// A is a public member of an array type provided by a dependency. This would
	// dependency bleeding but this is a main package.
	A [3]a.Type
	// B is an array of a standard type. Nothing to see here.
	B [3]int
}
