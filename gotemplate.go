package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Embed the template files (excluding .gitignore, gotemplate.go, and cmd/*).
// We now embed goreleaser.yaml (without a dot) along with other desired files.
//
//go:embed goreleaser.yaml README.md go.mod go.sum Makefile LICENSE
var templateFiles embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gotemplate <projectname>")
		os.Exit(1)
	}
	projectName := os.Args[1]

	// Create the project root directory.
	if err := os.Mkdir(projectName, 0755); err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		os.Exit(1)
	}

	// Walk through the embedded template files and write them out.
	err := fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Replace "gotemplate" with the new project name in file names.
		newPath := strings.ReplaceAll(path, "gotemplate", projectName)
		destPath := filepath.Join(projectName, newPath)

		// If the file is goreleaser.yaml, change its name to .goreleaser.yaml in the destination.
		if filepath.Base(destPath) == "goreleaser.yaml" {
			destPath = filepath.Join(filepath.Dir(destPath), ".goreleaser.yaml")
		}

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := templateFiles.ReadFile(path)
		if err != nil {
			return err
		}

		var content string
		// For .goreleaser.yaml, perform the special substitution first.
		if filepath.Base(destPath) == ".goreleaser.yaml" {
			content = strings.ReplaceAll(string(data), "./gotemplate.go", "./cmd/main.go")
			content = strings.ReplaceAll(content, "gotemplate", projectName)
		} else {
			content = strings.ReplaceAll(string(data), "gotemplate", projectName)
			// For Makefile, remove " --config goreleaser.yaml"
			if filepath.Base(destPath) == "Makefile" {
				content = strings.ReplaceAll(content, " --config goreleaser.yaml", "")
			}
		}

		return os.WriteFile(destPath, []byte(content), 0644)
	})
	if err != nil {
		fmt.Printf("Error copying template files: %v\n", err)
		os.Exit(1)
	}

	// Write a new mypackage.go with an empty package declaration.
	mypackagePath := filepath.Join(projectName, "mypackage.go")
	if err := os.WriteFile(mypackagePath, []byte("package mypackage\n"), 0644); err != nil {
		fmt.Printf("Error writing mypackage.go: %v\n", err)
		os.Exit(1)
	}

	// Create the cmd directory and write a new main.go with an empty main() function.
	cmdDir := filepath.Join(projectName, "cmd")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		fmt.Printf("Error creating cmd directory: %v\n", err)
		os.Exit(1)
	}

	mainGoPath := filepath.Join(cmdDir, "main.go")
	mainContent := `package main

func main() {
}
`
	if err := os.WriteFile(mainGoPath, []byte(mainContent), 0644); err != nil {
		fmt.Printf("Error writing cmd/main.go: %v\n", err)
		os.Exit(1)
	}

	// Manually write the .gitignore file with the specified contents.
	gitignorePath := filepath.Join(projectName, ".gitignore")
	gitignoreContent := `.DS_Store
dist
`
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		fmt.Printf("Error writing .gitignore: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project %s created successfully.\n", projectName)
}
