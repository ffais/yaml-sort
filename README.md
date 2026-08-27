# YAML Sort

[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](https://github.com/ffais/yaml-sort/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ffais/yaml-sort)](https://goreportcard.com/report/github.com/ffais/yaml-sort)
[![Go Reference](https://pkg.go.dev/badge/github.com/ffais/yaml-sort.svg)](https://pkg.go.dev/github.com/ffais/yaml-sort)

A command-line tool to rewrite YAML files with deterministic key ordering while retaining comments and anchors.

## Features

- **Preserves comments** - Maintains all inline and standalone comments
- **Customizable sorting** - Choose between alphabetical or custom key order
- **Configurable formatting** - Controls indentation and spacing between top-level keys
- **Multi-document support** - Handles YAML files with multiple documents (`---` separators)
- **Recursive search** - Sorts every file with a given name under a directory

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
yaml-sort sort --input-file input.yaml --output-file sorted.yaml
```

Sort with custom key order:
```bash
yaml-sort sort -i input.yaml -o sorted.yaml -c "name,version,dependencies"
```

Sort in-place (modify file directly):
```bash
yaml-sort sort --in-place --input-file input.yaml
yaml-sort sort --in-place config/one.yaml config/two.yaml
```

Recursively sort files with the same name:
```bash
yaml-sort sort --search-dir ./services --input-file values.yaml
```

## Options

```
  -i, --input-file string    Input YAML file
  -o, --output-file string   Output YAML file
  -w, --in-place             Replace input files with sorted content
  -c, --custom-sort strings  Comma-separated keys to place first
  -r, --reverse              Reverse sort order
  -l, --sort-list            Sort sequence values
  -s, --space-top-key        Add spacing between top-level keys (default true)
  -t, --indent int           Indentation width (default 2)
  -d, --search-dir string    Recursively find files named by --input-file
      --config string        Read options from a configuration file
  -h, --help                 Show help
```

Run `yaml-sort version` to display build information. Output must be provided with
`--output-file` unless `--in-place` or `--search-dir` is used.

## Examples

1. Sort a Kubernetes manifest:
```bash
yaml-sort sort -i deployment.yaml -o sorted-deployment.yaml
```

2. Sort with specific key priority:
```bash
yaml-sort sort -i config.yaml -o sorted-config.yaml -c "apiVersion,kind,metadata,spec"
```

3. Sort all `values.yaml` files under a directory:
```bash
yaml-sort sort -d ./clusters -i values.yaml
```

## Contributing

Pull requests are welcome! Please ensure:
- Tests pass (`go test ./...`)
- Code is formatted (`go fmt`)
- New features include appropriate tests

## License

Apache 2.0 - See [LICENSE](LICENSE) file for details.
```
