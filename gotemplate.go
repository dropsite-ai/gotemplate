package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
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
		fmt.Println("Usage:")
		fmt.Println("  gotemplate <projectname>  - Create a new Go project from the template.")
		fmt.Println("  gotemplate -commit        - Commit known template files separately.")
		os.Exit(1)
	}

	// If the user just wants to commit the template files in the current directory.
	if os.Args[1] == "-commit" {
		if err := commitTemplateFiles(); err != nil {
			fmt.Printf("Error committing template files: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Template files committed successfully.")
		return
	}

	// Otherwise, create a new project with the provided name.
	projectName := os.Args[1]
	if err := createProject(projectName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Project %s created successfully.\n", projectName)
}

func commitTemplateFiles() error {
	// Define the file groups and their associated commit messages.
	groups := []struct {
		files   []string
		message string
	}{
		{
			files:   []string{"Makefile"},
			message: "Makefile",
		},
		{
			files:   []string{"README.md"},
			message: "Project README",
		},
		{
			files:   []string{"go.mod", "go.sum"},
			message: "Go modules",
		},
		{
			files:   []string{".gitignore"},
			message: "Git ignore paths",
		},
		{
			files:   []string{".goreleaser.yaml"},
			message: "Goreleaser configuration",
		},
	}

	for _, group := range groups {
		// `git add` the files in this group.
		addArgs := append([]string{"add"}, group.files...)
		if err := runGitCommand(addArgs...); err != nil {
			return err
		}

		// Commit them with the specified commit message.
		commitArgs := []string{"commit", "-m", group.message}
		if err := runGitCommand(commitArgs...); err != nil {
			return err
		}
	}
	return nil
}

func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createProject(projectName string) error {
	// Create the project root directory.
	if err := os.Mkdir(projectName, 0755); err != nil {
		return fmt.Errorf("creating project directory: %v", err)
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
		return fmt.Errorf("copying template files: %v", err)
	}

	// Write a new file using the project name as the package declaration.
	packageFilePath := filepath.Join(projectName, projectName+".go")
	packageContent := fmt.Sprintf("package %s\n", projectName)
	if err := os.WriteFile(packageFilePath, []byte(packageContent), 0644); err != nil {
		return fmt.Errorf("writing %s.go: %v", projectName, err)
	}

	// Create the cmd directory and write a new main.go with an empty main() function.
	cmdDir := filepath.Join(projectName, "cmd")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		return fmt.Errorf("creating cmd directory: %v", err)
	}

	mainGoPath := filepath.Join(cmdDir, "main.go")
	mainContent := `package main

func main() {
}
`
	if err := os.WriteFile(mainGoPath, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("writing cmd/main.go: %v", err)
	}

	// Manually write the .gitignore file with the specified contents.
	gitignorePath := filepath.Join(projectName, ".gitignore")
	gitignoreContent := `.DS_Store
dist
`
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("writing .gitignore: %v", err)
	}

	return nil
}
