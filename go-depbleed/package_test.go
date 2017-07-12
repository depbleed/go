package depbleed

import (
	"fmt"
	"go/token"
	"go/types"
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

func TestGetPackageInfo(t *testing.T) {
	info, err := GetPackageInfo("github.com/depbleed/go/go-depbleed")

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	expected := "depbleed"

	if info.Package.Name() != expected {
		t.Errorf("expected \"%s\", got \"%s\"", expected, info.Package.Name())
	}
}

func TestGetPackageInfoNonExistingPackage(t *testing.T) {
	_, err := GetPackageInfo("github.com/depbleed/go/go-depbleed/nonexisting")

	if err == nil {
		t.Error("expected an error but didn't get one")
	}
}

func TestGetTypePackagePathBasicType(t *testing.T) {
	expected := ""
	v := &types.Basic{}
	path := GetTypePackagePath(v)

	if path != expected {
		t.Errorf("expected \"%s\" got \"%s\"", expected, path)
	}
}

func TestGetTypePackagePathNamedType(t *testing.T) {
	expected := "foo/bar"
	pkg := types.NewPackage("foo/bar", "bar")
	typename := types.NewTypeName(token.NoPos, pkg, "MyType", types.NewStruct(nil, nil))
	v := types.NewNamed(typename, &types.Basic{}, nil)
	path := GetTypePackagePath(v)

	if path != expected {
		t.Errorf("expected \"%s\" got \"%s\"", expected, path)
	}
}

func TestGetTypeShortName(t *testing.T) {
	expected := "bar.MyType"
	pkg := types.NewPackage("foo/bar", "bar")
	typename := types.NewTypeName(token.NoPos, pkg, "MyType", types.NewStruct(nil, nil))
	v := types.NewNamed(typename, &types.Basic{}, nil)
	shortName := GetTypeShortName(v)

	if shortName != expected {
		t.Errorf("expected \"%s\" got \"%s\"", expected, shortName)
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

func TestIsSubPackage(t *testing.T) {
	rootPackage := "github.com/depbleed/go/examples/exstruct"
	testCases := []struct {
		Package  string
		Expected bool
	}{
		{
			Package:  "github.com/depbleed/go/examples/exstruct/vendor/",
			Expected: false,
		},
		{
			Package:  "github.com/depbleed/go/examples/exstruct",
			Expected: true,
		},
		{
			Package:  "github.com/depbleed/go/examples/exstruct/sub",
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Package, func(t *testing.T) {
			value := IsSubPackage(testCase.Package, rootPackage)

			if value != testCase.Expected {
				t.Errorf("expected %t but got %t", testCase.Expected, value)
			}
		})
	}
}
