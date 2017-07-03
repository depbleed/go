package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/loader"

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

		absPath, err := filepath.Abs(path)

		if err != nil {
			return fmt.Errorf("could not understand path \"%s\": %s", path, err)
		}

		gopath := os.Getenv("GOPATH")

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
			var config loader.Config
			config.Import(packagePath)
			program, err := config.Load()

			if err != nil {
				return err
			}

			info := program.Package(packagePath)

			// TODO: Remove this. For debugging purposes only.
			for name, def := range info.Defs {
				if def != nil {
					fmt.Printf("%s: %s\n", name, def.String())
				} else {
					fmt.Println("name", name)
				}
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
