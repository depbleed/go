package exstruct

import "a"
import "net/http"

// MyType is an exposed type.
type MyType struct {
	// A is a public member of a type provided by a dependency. This is
	// dependency bleeding.
	A a.Type
	//B is a public member of a type provided by Go; nothing wrong here.
	B string
	//C is a public member of a type by this package; nothing wrong here.
	C MyOtherType
	//D is a public member of a type provided by a Go package; nothing wrong here.
	D http.Client
}

//MyOtherType is another type exposed here
type MyOtherType struct {
}
