package main

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	depbleed "github.com/depbleed/go/go-depbleed"
	"github.com/spf13/cobra"
)

// LintingError indicates a linting error occured.
type LintingError struct{}

func (LintingError) Error() string { return "linting error" }

var (
	noFail bool
)

var rootCmd = cobra.Command{
	Use:   "depbleed [path/package]",
	Short: "A Go linter that reports dependency bleeding",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("too many arguments")
		}

		path := "."

		if len(args) == 1 {
			path = args[0]
		}

		wd, err := os.Getwd()

		if err != nil {
			return fmt.Errorf("failed to get working directory: %s", err)
		}

		gopath := build.Default.GOPATH

		cmd.SilenceUsage = true

		packagePaths, err := depbleed.GetPackagePaths(gopath, path)

		if err != nil {
			return fmt.Errorf("could not get package paths: %s", err)
		}

		failed := false

		for _, packagePath := range packagePaths {
			packageInfo, err := depbleed.GetPackageInfo(packagePath)

			if err != nil {
				return err
			}

			leaks := packageInfo.Leaks()

			for _, leak := range leaks {
				relPath, err := filepath.Rel(wd, leak.Position.Filename)

				if err != nil {
					relPath = leak.Position.Filename
				}

				fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", relPath, leak.Position.Line, leak.Position.Column, leak)
			}

			if len(leaks) != 0 {
				failed = true
			}
		}

		if !noFail && failed {
			return LintingError{}
		}

		return nil
	},
	SilenceErrors: true,
}

func init() {
	rootCmd.Flags().BoolVar(&noFail, "no-fail", false, "Don't fail on errors")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		switch err.(type) {
		case LintingError:
		default:
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		os.Exit(127)
	}
}
