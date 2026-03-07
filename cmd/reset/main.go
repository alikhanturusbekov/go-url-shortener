package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alikhanturusbekov/go-url-shortener/internal/tools/resetgen"
)

// main scans the project tree, finds marked structs, and generates Reset methods
func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			switch d.Name() {
			case ".git", ".idea", "vendor":
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".go" || filepath.Base(path) == "reset.gen.go" {
			return nil
		}

		structs, err := resetgen.ParseFile(path)
		if err != nil {
			return err
		}

		if len(structs) == 0 {
			return nil
		}

		pkg := structs[0].Package
		dir := structs[0].Dir

		code := resetgen.Generate(pkg, structs)
		return resetgen.WriteFile(dir, code)
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
