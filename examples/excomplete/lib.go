package exstruct

import (
	"a"
	"net/http"
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
type G struct {
	a.Interface
}

//H is standard type chan
var H chan a.Int

//I is a standard type pointer
var I *a.Time

//J is an exported standard type
var J a.Int = 1

// K has an exported map key.
var K map[a.Int]string

// L has an exported map value.
var L map[string]a.Int

// M has an exported item type.
var M []a.Int

// N has an exported item type.
var N [3]a.Int

// O exports a standard type
var O http.Client
