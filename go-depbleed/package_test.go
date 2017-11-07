package depbleed

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetPackagePaths(t *testing.T) {
	fixturesGoPath, _ := filepath.Abs("./fixtures/gopath")
	testCases := []struct {
		Gopath   string
		Path     string
		Expected []string
	}{
		{
			Gopath:   "/tmp",
			Path:     "/tmp/src/foo/bar",
			Expected: []string{"foo/bar"},
		},
		{
			Gopath:   "./tmp",
			Path:     "/tmp/foo/bar",
			Expected: nil,
		},
		{
			Gopath:   "/tmp",
			Path:     "/tmp/foo/bar",
			Expected: nil,
		},
		{
			Gopath:   "/tmp2",
			Path:     "/tmp/src/foo/bar",
			Expected: nil,
		},
		{
			Gopath:   "/tmp",
			Path:     "foo",
			Expected: []string{"foo"},
		},
		{
			Gopath:   fixturesGoPath,
			Path:     "./fixtures/gopath/src/foo/...",
			Expected: []string{"foo", "foo/bar"},
		},
		{
			Gopath:   fixturesGoPath,
			Path:     "./fixtures/gopath/src/foo/unexisting/...",
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s-%s", testCase.Gopath, testCase.Path), func(t *testing.T) {
			gopath := filepath.FromSlash(testCase.Gopath)
			path := filepath.FromSlash(testCase.Path)

			values, err := GetPackagePaths(gopath, path)

			if len(testCase.Expected) == 0 {
				if err == nil {
					t.Errorf("expected an error but got: %v", values)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %s", err)
				}

				if !reflect.DeepEqual(values, testCase.Expected) {
					t.Errorf("expected \"%v\" but got \"%v\"", testCase.Expected, values)
				}
			}
		})
	}
}

func TestUseVCSRootOption(t *testing.T) {
	// Coverage only.
	option := UseVCSRootOption("my/go/path")

	if option == nil {
		t.Fatal("expected not nil")
	}
}

func TestGetPackageInfo(t *testing.T) {
	info, err := GetPackageInfo("github.com/depbleed/go/go-depbleed")

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	expected := "depbleed"

	if info.Package.Name() != expected {
		t.Errorf("expected \"%s\", got \"%s\"", expected, info.Package.Name())
	}
}

type failOption struct{}

func (failOption) apply(*PackageInfo) error { return errors.New("fail") }

func TestGetPackageInfoOptionError(t *testing.T) {
	_, err := GetPackageInfo("github.com/depbleed/go/go-depbleed", failOption{})

	if err == nil {
		t.Error("expected an error but didn't get one")
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

func TestGetRoot(t *testing.T) {
	info := PackageInfo{
		Package: &types.Package{},
		VCSRoot: "",
	}

	expected := ""
	value := info.GetRoot()

	if value != expected {
		t.Errorf("expected\"%s\", got: \"%s\"", expected, value)
	}
}

func TestGetRootWithVCS(t *testing.T) {
	info := PackageInfo{
		Package: &types.Package{},
		VCSRoot: "foo/bar",
	}

	expected := "foo/bar"
	value := info.GetRoot()

	if value != expected {
		t.Errorf("expected\"%s\", got: \"%s\"", expected, value)
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
