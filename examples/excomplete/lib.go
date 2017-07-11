package exstruct

import (
	"a"
)

//A is an exported standard type
const A a.Int = 1

//B is an exported vendor type
var B = a.Struct{}

//C is an exported struct that exports a standard type
type C struct {
	D a.Time
}

//E is an exported interface that mixes vendor & standard type
type E interface {
	F(a.Int) a.Time
}

//G exports a standard type
type G interface {
	a.Interface
}

//H is standard type chan
var H chan a.Int

//I is a standard type pointer
var I *a.Time

//J is an exported standard type
var J a.Int = 1
