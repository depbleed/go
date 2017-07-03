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

func TestIsVendor(t *testing.T) {

	testCases := []struct {
		Root     string
		Class    string
		Expected bool
	}{
		{
			Root:     "github.com/depbleed/go/examples/exstruct",
			Class:    "github.com/depbleed/go/examples/exstruct/vendor/a.Type",
			Expected: true,
		},
		{
			Root:     "github.com/depbleed/go/examples/exstruct",
			Class:    "github.com/depbleed/go/examples/exstruct.MyOtherType",
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s-%s", testCase.Root, testCase.Class), func(t *testing.T) {

			value := isVendor(testCase.Root, testCase.Class)

			if value != testCase.Expected {
				t.Errorf("expected \"%t\" but got \"%t\"", testCase.Expected, value)
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
