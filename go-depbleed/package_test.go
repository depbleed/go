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
