package depbleed

import (
	"fmt"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"
)

func isFilePath(path string) bool {
	return filepath.IsAbs(path) || strings.HasPrefix(path, ".")
}

func isVendor(path string) bool {
	return path == "vendor"
}

func isHidden(path string) bool {
	return strings.HasPrefix(path, ".")
}

func goPackagesWalkFunc(gopath string, packages map[string]bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if isVendor(info.Name()) {
				return filepath.SkipDir
			}

			if isHidden(info.Name()) {
				return filepath.SkipDir
			}
		} else if filepath.Ext(path) == ".go" {
			packagePath, _ := filepath.Rel(filepath.Join(gopath, "src"), filepath.Dir(path))
			packages[packagePath] = true
		}

		return nil
	}
}

func scanGoPackages(gopath string, path string) (result []string, err error) {
	packages := make(map[string]bool)
	err = filepath.Walk(path, goPackagesWalkFunc(gopath, packages))

	if err != nil {
		return nil, fmt.Errorf("unable to walk through \"%s\": %s", path, err)
	}

	for path := range packages {
		result = append(result, path)
	}

	sort.Strings(result)

	return
}

// GetPackagePaths returns the package paths for the packages matching the
// specified `path`.
//
// If `path` is a Go package path, is it returned as-is. This is a convenience.
//
// If `path` is either an absolute file path, or starts with a dot, the
// specified `gopath` is used to determine the package paths.
//
// If the file `path` is not in `gopath` or if it's relative position to the
// `gopath` can't be determined, an error is returned.
//
// If `path` is a filepath and ends with ..., subpackages are also looked for
// recursively.
//
// GetPackagePaths does not check for the package existence and will hapilly
// return a package path for a non-existing package.
func GetPackagePaths(gopath string, path string) ([]string, error) {
	if isFilePath(path) {
		path, _ = filepath.Abs(path)
		packagePath, err := filepath.Rel(filepath.Join(gopath, "src"), path)

		if err != nil {
			return nil, fmt.Errorf("cannot determine if path \"%s\" is in GOPATH (%s): %s", path, gopath, err)
		}

		if strings.HasPrefix(packagePath, "..") {
			return nil, fmt.Errorf("path \"%s\" is not in GOPATH (%s)", path, gopath)
		}

		if strings.HasSuffix(packagePath, "...") {
			dir := filepath.Dir(path)

			return scanGoPackages(gopath, dir)
		}

		return []string{filepath.ToSlash(packagePath)}, nil
	}

	return []string{path}, nil
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

// IsMain checks whether the package is a main package.
func (i PackageInfo) IsMain() bool {
	return i.Package.Name() == "main"
}

// Leaks returns the leaks in the package.
func (i PackageInfo) Leaks() (result Leaks) {
	if i.IsMain() {
		return
	}

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

	sort.Sort(result)

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
				return fmt.Errorf("function argument %s is an external type: %s", nameOrIndex(vars, j), err)
			}
		}

		vars = t.Results()

		for j := 0; j < vars.Len(); j++ {
			if err := i.CheckLeaks(vars.At(j).Type()); err != nil {
				return fmt.Errorf("function result %s is an external type: %s", nameOrIndex(vars, j), err)
			}
		}

		return nil
	case *types.Chan:
		if err := i.CheckLeaks(t.Elem()); err != nil {
			return fmt.Errorf("channel of external type: %s", err)
		}

		return nil
	case *types.Pointer:
		if err := i.CheckLeaks(t.Elem()); err != nil {
			return fmt.Errorf("pointer to external type: %s", err)
		}

		return nil
	case *types.Array:
		if err := i.CheckLeaks(t.Elem()); err != nil {
			return fmt.Errorf("array item is an external type: %s", err)
		}

		return nil
	case *types.Slice:
		if err := i.CheckLeaks(t.Elem()); err != nil {
			return fmt.Errorf("slice item is an external type: %s", err)
		}

		return nil
	case *types.Map:
		if err := i.CheckLeaks(t.Key()); err != nil {
			return fmt.Errorf("map key is an external type: %s", err)
		}

		if err := i.CheckLeaks(t.Elem()); err != nil {
			return fmt.Errorf("map value is an external type: %s", err)
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
		return fmt.Errorf("%s is a vendorized type from %s", GetTypeShortName(t), pkgPath)
	}

	return fmt.Errorf("%s is a global type from %s", GetTypeShortName(t), pkgPath)
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

// GetTypeShortName returns the short type representation for a given type.
func GetTypeShortName(t types.Type) string {
	parts := strings.Split(t.String(), "/")

	return parts[len(parts)-1]
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
