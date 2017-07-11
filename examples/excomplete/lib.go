package exstruct

import (
	"a"
)

const A a.Int = 1

var B a.Struct = a.Struct{}

type C struct {
	D a.Time
}

type E interface {
	F(a.Int) a.Time
}

type G interface {
	a.Interface
}

var H chan a.Int

var I *a.Time

var J a.Int = 1
