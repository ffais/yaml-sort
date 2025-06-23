package internal

import (
	"log"
	"os"
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
