package exchan

import "a"

// MyType is an exposed type.
type MyType struct {
	// A is a public member of a channel type provided by a dependency. This is
	// dependency bleeding.
	A chan a.Type
	// B is a channel of a standard type. Nothing to see here.
	B chan int
}
