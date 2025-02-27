# gotemplate

gotemplate for dropsite-ai.

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

```bash
gotemplate <name>
```

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