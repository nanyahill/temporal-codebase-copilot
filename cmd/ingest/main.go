package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
)

func main() {
	// Define flags
	repo := flag.String("repo", "", "path to the Temporal repo")
	limit := flag.Int("limit", 1, "number of files to process")

	// Parse flags
	flag.Parse()

	// Define fileset to keep track of file content
	fset := token.NewFileSet()

	// Basic validation
	if *repo == "" {
		log.Fatal("missing required --repo flag")
	}

	// Print values (for now, just to verify everything works)
	//fmt.Println("Repo path:", *repo)
	//fmt.Println("File limit:", *limit)

	count := 0
	err := filepath.WalkDir(*repo, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err // Stop if we hit a real problem
		}

		if d.IsDir() {
			if d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if count >= *limit {
			return filepath.SkipAll
		}
		extension := filepath.Ext(path)

		// Check if it's a Go file
		if extension == ".go" {
			fmt.Println("Found Go file:", path)

			// Parse the individual file
			// We pass 'nil' as the third argument because we want it to read from the 'path'
			node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil // Skip files that can't be parsed
			}

			// The package name is right here!
			packageName := node.Name.Name
			fmt.Printf("Package: %s\n", packageName)
			// Walk through the code "Tree"
			ast.Inspect(node, func(n ast.Node) bool {
				// Check: Is this node a Function Declaration?
				if fn, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function: %s\n", fn.Name.Name)
				}
				return true // Continue inspecting the rest of the file
			})
			count++
		}

		return nil // Keep going!
	})

	if err != nil {
		fmt.Println("Errors: ", err)
	}
}
