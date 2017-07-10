package depbleed

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestGetPackagePath(t *testing.T) {
	testCases := []struct {
		Gopath   string
		Path     string
		Expected string
	}{
		{
			Gopath:   "/tmp",
			Path:     "/tmp/src/foo/bar",
			Expected: "foo/bar",
		},
		{
			Gopath:   "./tmp",
			Path:     "/tmp/foo/bar",
			Expected: "",
		},
		{
			Gopath:   "/tmp",
			Path:     "/tmp/foo/bar",
			Expected: "",
		},
		{
			Gopath:   "/tmp2",
			Path:     "/tmp/src/foo/bar",
			Expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s-%s", testCase.Gopath, testCase.Path), func(t *testing.T) {
			gopath := filepath.FromSlash(testCase.Gopath)
			path := filepath.FromSlash(testCase.Path)

			value, err := GetPackagePath(gopath, path)

			if testCase.Expected == "" {
				if err == nil {
					t.Errorf("expected an error but got: %s", value)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %s", err)
				}

				if value != testCase.Expected {
					t.Errorf("expected \"%s\" but got \"%s\"", testCase.Expected, value)
				}
			}
		})
	}
}

func TestIsStandardPackage(t *testing.T) {
	testCases := []struct {
		Package  string
		Expected bool
	}{
		{
			Package:  "net",
			Expected: true,
		},
		{
			Package:  "net/http",
			Expected: true,
		},
		{
			Package:  "foo",
			Expected: false,
		},
		{
			Package:  "zzzzz",
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Package, func(t *testing.T) {
			result := IsStandardPackage(testCase.Package)

			if result != testCase.Expected {
				t.Errorf("expected %t but got %t", testCase.Expected, result)
			}
		})
	}
}

func TestIsVendorPackage(t *testing.T) {
	rootPackage := "github.com/depbleed/go/examples/exstruct"
	testCases := []struct {
		Package  string
		Expected bool
	}{
		{
			Package:  "github.com/depbleed/go/examples/exstruct/vendor/foo/bar",
			Expected: true,
		},
		{
			Package:  "github.com/depbleed/go/examples/exstruct",
			Expected: false,
		},
		{
			Package:  "github.com/depbleed/go/examples/exinterface/vendor/foo/bar",
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Package, func(t *testing.T) {
			value := IsVendorPackage(testCase.Package, rootPackage)

			if value != testCase.Expected {
				t.Errorf("expected %t but got %t", testCase.Expected, value)
			}
		})
	}
}

func TestIsSamePackage(t *testing.T) {

	testCases := []struct {
		Root     string
		Class    string
		Expected bool
	}{
		{
			Root:     "github.com/depbleed/go/examples/exstruct",
			Class:    "github.com/depbleed/go/examples/exstruct/vendor/a.Type",
			Expected: false,
		},
		{
			Root:     "github.com/depbleed/go/examples/exstruct",
			Class:    "github.com/depbleed/go/examples/exstruct.MyOtherType",
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s-%s", testCase.Root, testCase.Class), func(t *testing.T) {

			value := isSamePackage(testCase.Root, testCase.Class)

			if value != testCase.Expected {
				t.Errorf("expected \"%t\" but got \"%t\"", testCase.Expected, value)
			}

		})
	}
}

func TestIsNativeGo(t *testing.T) {

	testCases := []struct {
		Class    string
		Expected bool
	}{
		{
			Class:    "string",
			Expected: true,
		},
		{
			Class:    "net/http.Client",
			Expected: true,
		},
		{
			Class:    "github.com/depbleed/go/examples/exstruct.MyType",
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Class), func(t *testing.T) {

			value := isNativeGo(testCase.Class)

			if value != testCase.Expected {
				t.Errorf("expected \"%t\" but got \"%t\"", testCase.Expected, value)
			}

		})
	}
}

func TestIsLeaking(t *testing.T) {

	testCases := []struct {
		Root     string
		Class    string
		Expected bool
	}{
		{
			Root:     "github.com/depbleed/go/examples/exstruct",
			Class:    "string",
			Expected: false,
		},
		{
			Root:     "github.com/depbleed/go/examples/exstruct",
			Class:    "net/http.Client",
			Expected: false,
		},
		{
			Root:     "github.com/depbleed/go/examples/exstruct",
			Class:    "github.com/depbleed/go/examples/exstruct.MyType",
			Expected: false,
		},
		{
			Root:     "github.com/depbleed/go/examples/exstruct",
			Class:    "github.com/depbleed/go/examples/exstruct/vendor/a.Type",
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Class), func(t *testing.T) {

			value := IsLeaking(testCase.Root, testCase.Class)

			if value != testCase.Expected {
				t.Errorf("expected \"%t\" but got \"%t\"", testCase.Expected, value)
			}

		})
	}
}
