[![Build Status](https://travis-ci.org/depbleed/go.png)](https://travis-ci.org/depbleed/go)
[![GoDoc](https://godoc.org/github.com/depbleed/go/go-depbleed?status.png)](https://godoc.org/github.com/depbleed/go/go-depbleed)
[![codecov](https://codecov.io/gh/depbleed/go/branch/master/graph/badge.svg)](https://codecov.io/gh/depbleed/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/depbleed/go)](https://goreportcard.com/report/github.com/depbleed/go)

# Depbleed

A Go linter that detects dependency-bleeding in Go packages.

## Rationale

Go encourages a vendoring model for dependencies where projects have to copy -
or fetch upon build - their dependencies and put them in a `vendor` folder. In
Go, packages are uniquely named according to their relative path to the
`GOPATH`. This causes situations where two libraries that use types from the
same dependency but in different `vendor` folders have identical yet
incompatible definitions for those types.

Consider the following layout in your library:

```
yourlib/
  vendor/
    lib-a/
      vendor/
        lib-foo/
          foo.go
    lib-b/
      vendor/
        lib-foo/
          foo.go
```

If `foo.go` contains a type `Foo` that both `lib-a` and `lib-b` **expose**, an
instance of `Foo` from `lib-a` won't be usable as-is by `lib-b`. Go will
complain that those are incompatible types (and rightly so).

While this may seem surprising, the real problem is not with Go's behavior but
in the fact that both `lib-a` and `lib-b` **bleed** a type from their
dependencies instead of providing an abstraction. Ideally, their dependencies
should remain an implementation detail.

The only exception, of course, is for standard types and libraries, which are
by definition the same for all libraries in your ecosystem.

Depbleed aims to detect dependency bleeding in your libraries, so that they
don't cause the same problems for your users.

Check out the [examples](examples/) for real-cases of dependency bleeding.
