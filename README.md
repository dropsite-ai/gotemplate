# gotemplate

Go template for Dropsite.

## Installation

### Homebrew (macOS or Linux)
```bash
brew tap dropsite-ai/homebrew-tap
brew install gotemplate
```

### Go
```bash
go get github.com/dropsite-ai/gotemplate
```

### Direct Download (Prebuilt Binaries)

1. Grab the latest binary from [GitHub Releases](https://github.com/dropsite-ai/gotemplate/releases).  
2. Extract and run the `gotemplate` executable.

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/dropsite-ai/gotemplate.git
   cd gotemplate
   ```
2. Build:
   ```bash
   go build -o gotemplate cmd/main.go
   ```

## Usage

### Create a New Project

To create a new Go project using the predefined template:
```bash
gotemplate <projectname>
```
This will generate a new directory named `<projectname>` containing common project files such as:
- `Makefile`
- `README.md`
- `go.mod`, `go.sum`
- `.goreleaser.yaml`
- A basic `cmd/main.go`
- A `.gitignore` file

### Commit Template Files Separately

If you are inside a Git repository and want to quickly commit template files in logical groups, use:
```bash
gotemplate -commit
```

This ensures cleaner commit history when initializing a project.

## Test

```bash
make test
```

## Release

```bash
make release
```

## License

[MIT License](LICENSE)
```
