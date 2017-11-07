package depbleed

import (
	"os"
	"testing"
)

func TestExamples(t *testing.T) {
	testCases := []struct {
		PackagePath      string
		LeaksCount       int
		UseVCSLeaksCount int
	}{
		{
			PackagePath:      "github.com/depbleed/go/examples/exinterface",
			LeaksCount:       2,
			UseVCSLeaksCount: 2,
		},
		{
			PackagePath:      "github.com/depbleed/go/examples/exstruct",
			LeaksCount:       1,
			UseVCSLeaksCount: 1,
		},
		{
			PackagePath:      "github.com/depbleed/go/examples/exmap",
			LeaksCount:       2,
			UseVCSLeaksCount: 2,
		},
		{
			PackagePath:      "github.com/depbleed/go/examples/exslice",
			LeaksCount:       1,
			UseVCSLeaksCount: 1,
		},
		{
			PackagePath:      "github.com/depbleed/go/examples/exarray",
			LeaksCount:       1,
			UseVCSLeaksCount: 1,
		},
		{
			PackagePath:      "github.com/depbleed/go/examples/expointer",
			LeaksCount:       1,
			UseVCSLeaksCount: 1,
		},
		{
			PackagePath:      "github.com/depbleed/go/examples/exchan",
			LeaksCount:       1,
			UseVCSLeaksCount: 1,
		},
		{
			PackagePath:      "github.com/depbleed/go/examples/exvcs",
			LeaksCount:       1,
			UseVCSLeaksCount: 0,
		},
		{
			PackagePath:      "github.com/depbleed/go/examples/excomplete",
			LeaksCount:       13,
			UseVCSLeaksCount: 12,
		},
		{
			PackagePath: "github.com/depbleed/go/examples/exmain",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.PackagePath, func(t *testing.T) {
			packageInfo, err := GetPackageInfo(testCase.PackagePath)

			if err != nil {
				t.Fatalf("expected no error but got: %s", err)
			}

			leaks := packageInfo.Leaks()

			if len(leaks) != testCase.LeaksCount {
				t.Errorf("expected %d leak(s) got %d", testCase.LeaksCount, len(leaks))
			}
		})
		t.Run(testCase.PackagePath+" use-vcs", func(t *testing.T) {
			gopath := os.Getenv("GOPATH")
			packageInfo, err := GetPackageInfo(testCase.PackagePath, UseVCSRootOption(gopath))

			if err != nil {
				t.Fatalf("expected no error but got: %s", err)
			}

			leaks := packageInfo.Leaks()

			if len(leaks) != testCase.UseVCSLeaksCount {
				t.Errorf("expected %d leak(s) got %d", testCase.UseVCSLeaksCount, len(leaks))
			}
		})
	}
}
