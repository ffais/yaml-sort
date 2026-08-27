package internal

import (
	"os"
	"path/filepath"
	"testing"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
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

func TestWriteToFileReplacesExistingFileAtomically(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.yaml")
	if err := os.WriteFile(outputPath, []byte("old: content\n"), 0600); err != nil {
		t.Fatal(err)
	}

	var document yaml.Node
	if err := yaml.Unmarshal([]byte("new: content\n"), &document); err != nil {
		t.Fatal(err)
	}
	WriteToFile(outputPath, []*yaml.Node{&document}, Config{Indent: 2})

	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "new: content\n" {
		t.Fatalf("got output %q, want %q", output, "new: content\n")
	}
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("got mode %o, want 600", info.Mode().Perm())
	}
	temporaryFiles, err := filepath.Glob(filepath.Join(tempDir, ".output.yaml.tmp-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(temporaryFiles) != 0 {
		t.Fatalf("temporary files were not cleaned up: %v", temporaryFiles)
	}
}