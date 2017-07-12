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

var rootCmd = cobra.Command{
	Use:   "depbleed [path]",
	Short: "A Go package for dependency bleeding",
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

		absPath, err := filepath.Abs(path)

		if err != nil {
			return fmt.Errorf("could not understand path \"%s\": %s", path, err)
		}

		gopath := build.Default.GOPATH

		cmd.SilenceUsage = true

		// TODO: The current code does not handle the `./...` wildcard form but
		// it should. Which means we must also handle a list of package paths,
		// not just a single one.
		var packagePaths []string

		packagePath, err := depbleed.GetPackagePath(gopath, absPath)

		if err != nil {
			return fmt.Errorf("could not get package path: %s", err)
		}

		packagePaths = append(packagePaths, packagePath)

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

			if len(leaks) == 0 {
				fmt.Fprintf(os.Stdout, "No leak detected for package %s\n", packagePath)
			}
		}

		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
