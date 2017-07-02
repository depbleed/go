package type_bleed_interface

import (
	"fake.com/lib-a"
)

type MyInterface interface {
	GetA() lib_a.LibAType
}
