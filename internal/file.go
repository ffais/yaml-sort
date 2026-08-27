package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	yaml "go.yaml.in/yaml/v3"
)

func ParseYaml(filePath string) ([]*yaml.Node, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("read YAML file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var documents []*yaml.Node
	for {
		var document yaml.Node
		err = decoder.Decode(&document)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decode YAML document: %w", err)
		}
		documents = append(documents, &document)
	}
	return documents, nil
}

func WriteToFile(filePath string, documents []*yaml.Node, cfg Config) error {
	// Marshal back to YAML with comments
	var output strings.Builder
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(cfg.Indent)
	for _, document := range documents {
		if err := encoder.Encode(document); err != nil {
			return fmt.Errorf("encode YAML document: %w", err)
		}
	}
	if err := writeFileAtomically(filePath, []byte(output.String())); err != nil {
		return fmt.Errorf("write YAML file: %w", err)
	}
	return nil
}

func writeFileAtomically(filePath string, data []byte) error {
	mode := os.FileMode(0644)
	if info, err := os.Stat(filePath); err == nil {
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	tempFile, err := os.CreateTemp(filepath.Dir(filePath), "."+filepath.Base(filePath)+".tmp-*")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if err := tempFile.Chmod(mode); err != nil {
		tempFile.Close()
		return err
	}
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return err
	}
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	return os.Rename(tempPath, filePath)
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
