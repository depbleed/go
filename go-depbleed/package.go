package depbleed

import (
	"fmt"
	"path/filepath"
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

func isVendor(rootPackage string, class string) bool {

	if strings.HasPrefix(class, rootPackage) && strings.Contains(class, "/vendor/") {
		return true
	}

	return false
}
