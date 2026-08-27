package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAndWriteMultipleDocuments(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.yaml")
	outputPath := filepath.Join(tempDir, "output.yaml")
	input := []byte("b: 1\na: 2\n---\nd: 3\nc: 4\n")
	if err := os.WriteFile(inputPath, input, 0600); err != nil {
		t.Fatal(err)
	}

	documents := ParseYaml(inputPath)
	if len(documents) != 2 {
		t.Fatalf("got %d documents, want 2", len(documents))
	}
	for _, document := range documents {
		SortYamlNodes(document, Config{})
	}
	WriteToFile(outputPath, documents, Config{Indent: 2})

	roundTripDocuments := ParseYaml(outputPath)
	if len(roundTripDocuments) != 2 {
		t.Fatalf("round trip produced %d documents, want 2", len(roundTripDocuments))
	}
}