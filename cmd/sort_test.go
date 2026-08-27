package cmd

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	internal "github.com/ffais/yaml-sort/internal"
)

func TestParallelProcessingReturnsWorkerErrors(t *testing.T) {
	worker := func(inputFile string, outputFile string, cfg internal.Config) error {
		return errors.New(inputFile)
	}

	err := parallelProcessing([]string{"first.yaml", "second.yaml"}, 2, worker)
	if err == nil {
		t.Fatal("parallelProcessing returned nil error")
	}
	for _, file := range []string{"first.yaml", "second.yaml"} {
		if !strings.Contains(err.Error(), file) {
			t.Errorf("error %q does not mention %q", err, file)
		}
	}
}

func TestInitConfigReturnsMissingFileError(t *testing.T) {
	originalCfgFile := cfgFile
	t.Cleanup(func() {
		cfgFile = originalCfgFile
	})
	cfgFile = filepath.Join(t.TempDir(), "missing.yaml")

	if err := initConfig(nil, nil); err == nil {
		t.Fatal("initConfig returned nil error for a missing config file")
	}
}
