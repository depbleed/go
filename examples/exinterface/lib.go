package exinterface

import (
	"a"
)

// MyInterface is an interface that does not do proper encapsulation.
type MyInterface interface {
	// GetA returns one of its dependencies type: this is dependency bleeding.
	GetA() a.Type
	// SetA takes one of its dependencies type: this is dependency bleeding.
	SetA(value a.Type)
	// OK is fine.
	OK(int) int
}
