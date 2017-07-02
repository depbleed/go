# Depbleed

A Go linter that detects dependency-bleeding in Go packages.

## Rationale

Go encourages a vendoring model for dependencies where projects have to copy -
or fetch upon build - their dependencies and put them in a `vendor` folders. In
Go, packages are uniquely named according to their relative path to the
`GOPATH`. This causes situations where two libraries that use types from the
same dependency but in different `vendor` folder have identical yet
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

Depbleed aims to detect such problems in your libraries, so that they don't
cause the same problems for your users.
