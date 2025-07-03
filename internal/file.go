package internal

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

func ParseYaml(filePath string, node *yaml.Node) {
	// Read the input file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	// Parse the YAML content
	err = yaml.Unmarshal(data, node)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
}

func WriteToFile(filePath string, node *yaml.Node, cfg Config) {
	// Marshal back to YAML with comments
	var output strings.Builder
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(cfg.Indent)
	err := encoder.Encode(node)
	if err != nil {
		log.Fatalf("Error marshaling YAML: %v", err)
	}
	// Write to output file
	err = os.WriteFile(filePath, []byte(output.String()), 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
}

func FindYamlFile(searchRoot string, file string) ([]string, error) {

	var yamlSearchRoot string
	var yamlFiles []string

	if path.IsAbs(searchRoot) {
		yamlSearchRoot = searchRoot
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("error getting working directory: %w", err)
		}

		yamlSearchRoot = filepath.Join(cwd, searchRoot)
	}

	err := filepath.WalkDir(yamlSearchRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		absolutePath, _ := filepath.Abs(path)
		// fmt.Printf("Visited: %s\n", path)

		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		if filepath.Base(path) == file {
			yamlFiles = append(yamlFiles, absolutePath)
		}
		return nil
	})
	return yamlFiles, err
}
