package depbleed

import "testing"

func TestExamples(t *testing.T) {
	testCases := []struct {
		PackagePath string
	}{
		{
			PackagePath: "github.com/depbleed/go/examples/exinterface",
		},
		{
			PackagePath: "github.com/depbleed/go/examples/exstruct",
		},
		{
			PackagePath: "github.com/depbleed/go/examples/exmap",
		},
		{
			PackagePath: "github.com/depbleed/go/examples/exslice",
		},
		{
			PackagePath: "github.com/depbleed/go/examples/exarray",
		},
		{
			PackagePath: "github.com/depbleed/go/examples/expointer",
		},
		{
			PackagePath: "github.com/depbleed/go/examples/exchan",
		},
		{
			PackagePath: "github.com/depbleed/go/examples/excomplete",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.PackagePath, func(t *testing.T) {
			packageInfo, err := GetPackageInfo(testCase.PackagePath)

			if err != nil {
				t.Errorf("expected no error but got: %s", err)
			}

			leaks := packageInfo.Leaks()

			if len(leaks) == 0 {
				t.Error("expected leaks")
			}
		})
	}
}
