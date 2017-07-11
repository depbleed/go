package depbleed

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
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
	Identifier *ast.Ident
	Object     types.Object
	Position   token.Position
}

// Leaks returns the leaks in the package.
func (i PackageInfo) Leaks() (result []Leak) {
	for identifier, obj := range i.Info.Defs {
		// Only exported types matter.
		if obj != nil && obj.Exported() {
			if i.IsLeaking(obj.Type()) {
				result = append(result, Leak{
					Identifier: identifier,
					Object:     obj,
					Position:   i.Fset.Position(obj.Pos()),
				})
			}
		}
	}

	return
}

// IsLeaking checks wheter a specified type is being leaked.
func (i PackageInfo) IsLeaking(t types.Type) bool {
	pkgPath := GetTypePackagePath(t)

	// Built-in type.
	if pkgPath == "" {
		return false
	}

	// Standard type.
	if IsStandardPackage(pkgPath) {
		return false
	}

	// Subpackages are ok.
	if IsSubPackage(pkgPath, i.Package.Path()) {
		return false
	}

	// Vendors are definitely leaking.
	if IsVendorPackage(pkgPath, i.Package.Path()) {
		return true
	}

	return false
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
