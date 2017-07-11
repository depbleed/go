package depbleed

import (
	"fmt"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"
)

// GetPackagePath returns the package path for the package at the specified
// location.
//
// If the `p` is not in `gopath`, an error is returned.
//
// GetPackagePath does not check for the path's existence and will hapilly
// return a package name for a non-existing package.
func GetPackagePath(gopath string, path string) (string, error) {
	packagePath, err := filepath.Rel(filepath.Join(gopath, "src"), path)

	if err != nil {
		return "", fmt.Errorf("cannot determine if path \"%s\" is in GOPATH (%s): %s", path, gopath, err)
	}

	if strings.HasPrefix(packagePath, "..") {
		return "", fmt.Errorf("path \"%s\" is not in GOPATH (%s)", path, gopath)
	}

	return filepath.ToSlash(packagePath), nil
}

// PackageInfo represents information about a package.
type PackageInfo struct {
	Package *types.Package
	Info    types.Info
	Fset    *token.FileSet
}

// GetPackageInfo returns information about the package at the specified
// location.
func GetPackageInfo(p string) (PackageInfo, error) {
	var config loader.Config
	config.Import(p)
	var nestedErr error

	config.TypeChecker.Error = func(err error) {
		nestedErr = err
	}

	program, err := config.Load()

	if err != nil {
		return PackageInfo{}, fmt.Errorf("%s: %s", err, nestedErr)
	}

	packageInfo := program.Package(p)

	return PackageInfo{
		Package: packageInfo.Pkg,
		Info:    packageInfo.Info,
		Fset:    config.Fset,
	}, nil
}

// Leak represents a leaking type.
type Leak struct {
	Object   types.Object
	Position token.Position
	err      error
}

func (l Leak) Error() string {
	return fmt.Sprintf("%s %s: %s", GetObjectKind(l.Object), l.Object.Name(), l.err)
}

// Leaks returns the leaks in the package.
func (i PackageInfo) Leaks() (result []Leak) {
	for _, obj := range i.Info.Defs {
		// Only exported types matter.
		if obj != nil && obj.Exported() {
			if err := i.CheckLeaks(obj.Type()); err != nil {
				result = append(result, Leak{
					Object:   obj,
					Position: i.Fset.Position(obj.Pos()),
					err:      err,
				})
			}
		}
	}

	return
}

// CheckLeaks checks wheter a specified type is being leaked.
func (i PackageInfo) CheckLeaks(t types.Type) error {
	switch t := t.(type) {
	case *types.Signature:
		vars := t.Params()

		nameOrIndex := func(t *types.Tuple, index int) string {
			name := t.At(index).Name()

			if name == "" {
				return strconv.Itoa(index)
			}

			return fmt.Sprintf("\"%s\"", name)
		}

		for j := 0; j < vars.Len(); j++ {
			if err := i.CheckLeaks(vars.At(j).Type()); err != nil {
				return fmt.Errorf("argument %s leaks internal type: %s", nameOrIndex(vars, j), err)
			}
		}

		vars = t.Results()

		for j := 0; j < vars.Len(); j++ {
			if err := i.CheckLeaks(vars.At(j).Type()); err != nil {
				return fmt.Errorf("result %s leaks internal type: %s", nameOrIndex(vars, j), err)
			}
		}

		return nil
	case *types.Chan:
		if err := i.CheckLeaks(t.Elem()); err != nil {
			return fmt.Errorf("channel of internal type: %s", err)
		}

		return nil
	case *types.Pointer:
		if err := i.CheckLeaks(t.Elem()); err != nil {
			return fmt.Errorf("pointer of internal type: %s", err)
		}

		return nil
	}

	pkgPath := GetTypePackagePath(t)

	// Built-in type.
	if pkgPath == "" {
		return nil
	}

	// Standard type.
	if IsStandardPackage(pkgPath) {
		return nil
	}

	// Subpackages are ok.
	if IsSubPackage(pkgPath, i.Package.Path()) {
		return nil
	}

	// Vendors are definitely leaking.
	if IsVendorPackage(pkgPath, i.Package.Path()) {
		return fmt.Errorf("%s is a vendorized type from %s", GetShortType(t), pkgPath)
	}

	return fmt.Errorf("%s is a non-local type from %s", GetShortType(t), pkgPath)
}

// GetTypePackagePath returns the package path for a given type.
//
// For built-in types (int, string, ...), an empty string is returned.
func GetTypePackagePath(t types.Type) string {
	parts := strings.Split(t.String(), ".")

	if len(parts) == 1 {
		return ""
	}

	return strings.Join(parts[:len(parts)-1], ".")
}

// GetShortType returns the short type representation for a given type.
func GetShortType(t types.Type) string {
	parts := strings.Split(t.String(), "/")

	return parts[len(parts)-1]
}

// GetObjectKind returns the kind of an object.
func GetObjectKind(o types.Object) (kind string) {
	switch o := o.(type) {
	case *types.Const:
		return fmt.Sprintf("%s constant", GetTypeKind(o.Type()))
	case *types.Var:
		return fmt.Sprintf("%s variable", GetTypeKind(o.Type()))
	}

	return GetTypeKind(o.Type())
}

// GetTypeKind returns the kind of a type.
func GetTypeKind(t types.Type) (kind string) {
	switch t := t.(type) {
	case *types.Named:
		return fmt.Sprintf("aliased %s", GetTypeKind(t.Underlying()))
	case *types.Struct:
		return "struct"
	case *types.Chan:
		return "channel type"
	case *types.Pointer:
		return "pointer type"
	case *types.Signature:
		return "function type"
	case *types.Basic:
		return "basic type"
	}
	return ""
}

// IsStandardPackage checks whether a given package is standard.
//
// Standard packages are provided with Go.
func IsStandardPackage(p string) bool {
	index := sort.SearchStrings(standardGoPackages, p)

	if index < len(standardGoPackages) {
		return standardGoPackages[index] == p
	}

	return false
}

// IsVendorPackage checks whether a given package is a vendor of the specified
// root package.
//
// A vendor is located in a `/vendor/` directory.
func IsVendorPackage(p string, rootPackage string) bool {
	return strings.HasPrefix(p, rootPackage) && strings.Contains(p, "/vendor/")
}

// IsSubPackage checks whether a given package is a subpackage of the specified
// root package.
//
// A package is always a subpackage of itself.
func IsSubPackage(p string, rootPackage string) bool {
	return strings.HasPrefix(p, rootPackage) && !strings.Contains(p, "/vendor/")
}
