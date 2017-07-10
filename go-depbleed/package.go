package depbleed

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// GetPackagePath returns the package path for the package at the specified
// location.
//
// If the `p` is not in `gopath`, an error is returned.
//
// GetPackagePath does not check for the path's existence and will hapilly
// return a package name for a non-existing package.
func GetPackagePath(gopath string, p string) (string, error) {
	packagePath, err := filepath.Rel(filepath.Join(gopath, "src"), p)

	if err != nil {
		return "", fmt.Errorf("cannot determine if path \"%s\" is in GOPATH (%s): %s", p, gopath, err)
	}

	if strings.HasPrefix(packagePath, "..") {
		return "", fmt.Errorf("path \"%s\" is not in GOPATH (%s)", p, gopath)
	}

	return filepath.ToSlash(packagePath), nil
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

//IsLeaking returns if `class` is being leaked by `rootPackage`
func IsLeaking(rootPackage string, class string) bool {

	return isVendor(rootPackage, class) && !isSamePackage(rootPackage, class) && !isNativeGo(class)
}

func isVendor(rootPackage string, class string) bool {

	if strings.HasPrefix(class, rootPackage) && strings.Contains(class, "/vendor/") {
		return true
	}

	return false
}

func isSamePackage(rootPackage string, class string) bool {
	if strings.HasPrefix(class, rootPackage) && !strings.Contains(class, "/vendor/") {
		return true
	}

	return false
}

func isNativeGo(class string) bool {
	//string, int and other base type do not
	//belong to a give package
	if !strings.Contains(class, "/") {
		return true
	}

	for _, goPackage := range standardGoPackages {

		if strings.HasPrefix(class, goPackage) {
			return true
		}
	}

	return false
}
