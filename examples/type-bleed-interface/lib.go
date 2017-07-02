package type_bleed_interface

import (
	"fake.com/lib-a"
)

// MyInterface is an interface that does not do proper encapsulation.
type MyInterface interface {
	// GetA returns one of its dependencies type: this is dependency bleeding.
	GetA() lib_a.LibAType
}
