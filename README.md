# YAML Sort

[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](https://github.com/ffais/yaml-sort/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ffais/yaml-sort)](https://goreportcard.com/report/github.com/ffais/yaml-sort)
[![Go Reference](https://pkg.go.dev/badge/github.com/ffais/yaml-sort.svg)](https://pkg.go.dev/github.com/ffais/yaml-sort)

A command-line tool to sort YAML files by keys while preserving comments, structure, and formatting.

## Features

- **Preserves comments** - Maintains all inline and standalone comments
- **Customizable sorting** - Choose between alphabetical or custom key order
- **Format preservation** - Keeps original indentation and style
- **Multi-document support** - Handles YAML files with multiple documents (`---` separators)
- **Strict mode** - Optionally validates YAML during processing

## Installation

### Using Go:
```bash
go install github.com/ffais/yaml-sort@latest
```

### Pre-built binaries:
Download from the [Releases page](https://github.com/ffais/yaml-sort/releases)

## Usage

Basic sorting:
```bash
yaml-sort -f input.yaml -o sorted.yaml
```

Sort with custom key order:
```bash
yaml-sort -f input.yaml -o sorted.yaml -k "name,version,dependencies"
```

Sort in-place (modify file directly):
```bash
yaml-sort -f input.yaml -i
```

## Options

```
  -f, --file string       Input YAML file
  -o, --output string     Output file (defaults to stdout if not specified)
  -i, --in-place          Modify input file directly
  -k, --keys string       Comma-separated list of keys for custom ordering
  -s, --strict            Enable strict YAML parsing
  -v, --version           Display version information
  -h, --help              Show help message
```

## Examples

1. Sort a Kubernetes manifest:
```bash
yaml-sort -f deployment.yaml -o sorted-deployment.yaml
```

2. Sort with specific key priority:
```bash
yaml-sort -f config.yaml -k "apiVersion,kind,metadata,spec"
```

3. Pipe input/output:
```bash
cat input.yaml | yaml-sort > sorted.yaml
```

## Contributing

Pull requests are welcome! Please ensure:
- Tests pass (`go test ./...`)
- Code is formatted (`go fmt`)
- New features include appropriate tests

## License

MIT - See [LICENSE](LICENSE) file for details.
```
