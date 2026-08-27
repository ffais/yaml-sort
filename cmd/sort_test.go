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

func TestValidateSortFlags(t *testing.T) {
	originalInputFile := InputFile
	originalOutputFile := OutputFile
	originalInPlace := InPlace
	originalCfg := Cfg
	t.Cleanup(func() {
		InputFile = originalInputFile
		OutputFile = originalOutputFile
		InPlace = originalInPlace
		Cfg = originalCfg
	})

	tests := []struct {
		name      string
		inputFile string
		outputFile string
		inPlace   bool
		searchDir string
		args      []string
		wantError bool
	}{
		{name: "input and output", inputFile: "in.yaml", outputFile: "out.yaml"},
		{name: "single file in place", inputFile: "in.yaml", inPlace: true},
		{name: "positional files in place", inPlace: true, args: []string{"in.yaml"}},
		{name: "search by file name", inputFile: "values.yaml", searchDir: "services"},
		{name: "missing input", outputFile: "out.yaml", wantError: true},
		{name: "missing output", inputFile: "in.yaml", wantError: true},
		{name: "positional files without in place", args: []string{"in.yaml"}, wantError: true},
		{name: "in place with output", inputFile: "in.yaml", outputFile: "out.yaml", inPlace: true, wantError: true},
		{name: "search without file name", searchDir: "services", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			InputFile = test.inputFile
			OutputFile = test.outputFile
			InPlace = test.inPlace
			Cfg = internal.Config{SearchDir: test.searchDir}

			err := validateSortFlags(nil, test.args)
			if (err != nil) != test.wantError {
				t.Fatalf("validateSortFlags() error = %v, wantError %v", err, test.wantError)
			}
		})
	}
}
