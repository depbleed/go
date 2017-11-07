package exchan

import "github.com/depbleed/go/examples/exchan"

// MyType is an exposed type.
type MyType struct {
	// A is a public member of a VCS-local type. This is theoretical bleeding
	// but in practice is fine.
	A exchan.MyType
}
