package type_bleed_struct

import (
	"fake.com/lib-a"
)

// MyType is an exposed type.
type MyType struct {
	// A is a public member of a type provided by a dependency. This is
	// dependency bleeding.
	A lib_a.LibAType
}
