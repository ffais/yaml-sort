package cmd

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

	internal "github.com/ffais/yaml-sort/internal"
	"github.com/spf13/cobra"
)

var OutputFile string
var InPlace bool

var sortCmd = &cobra.Command{
	Use:     "sort",
	Short:   "Yaml-Sort sorts content of YAML files alphabetically.",
	Long:    `Yaml-Sort sorts content of YAML files alphabetically preserving comments, anchor and with support for custom order.`,
	PreRunE: validateSortFlags,
	RunE:    sort,
}

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.Flags().StringVarP(&OutputFile, "output-file", "o", "", "The YAML file to output sorted content to.")
	sortCmd.Flags().BoolVarP(&InPlace, "in-place", "w", false, "Replace input files with sorted content")
}

func sort(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return parallelProcessing(args, runtime.NumCPU()*2, sortYamlFile)
	}
	if Cfg.SearchDir != "" {
		parallelism := runtime.NumCPU() * 2
		yamls, err := internal.FindYamlFile(Cfg.SearchDir, InputFile)
		if err != nil {
			return fmt.Errorf("search for YAML files: %w", err)
		}
		if len(yamls) == 0 {
			return fmt.Errorf("no files named %q found under %q", InputFile, Cfg.SearchDir)
		}
		return parallelProcessing(yamls, parallelism, sortYamlFile)
	}
	outputFile := OutputFile
	if InPlace {
		outputFile = InputFile
	}
	return sortYamlFile(InputFile, outputFile, Cfg)
}

func validateSortFlags(cmd *cobra.Command, args []string) error {
	if Cfg.SearchDir != "" {
		if len(args) > 0 {
			return errors.New("positional files cannot be combined with --search-dir")
		}
		if InputFile == "" {
			return errors.New("--input-file is required with --search-dir")
		}
		return nil
	}
	if len(args) > 0 {
		if !InPlace {
			return errors.New("positional files require --in-place")
		}
		if InputFile != "" || OutputFile != "" {
			return errors.New("positional files cannot be combined with --input-file or --output-file")
		}
		return nil
	}
	if InputFile == "" {
		return errors.New("--input-file is required")
	}
	if InPlace {
		if OutputFile != "" {
			return errors.New("--in-place cannot be combined with --output-file")
		}
		return nil
	}
	if OutputFile == "" {
		return errors.New("--output-file is required unless --in-place is set")
	}
	return nil
}

func sortYamlFile(inputFile string, outputFile string, cfg internal.Config) error {
	fmt.Println("Sorting yaml file", inputFile)
	documents, err := internal.ParseYaml(inputFile)
	if err != nil {
		return fmt.Errorf("sort %q: %w", inputFile, err)
	}
	for _, document := range documents {
		internal.SortYamlNodes(document, cfg)
		if Cfg.SpaceTopKey {
			internal.AddEmptyLinesBeforeTopLevelKeys(document)
		}
	}
	if err := internal.WriteToFile(outputFile, documents, cfg); err != nil {
		return fmt.Errorf("sort %q: %w", inputFile, err)
	}
	return nil
}

func parallelProcessing(files []string, parallelism int, fn func(inputFile string, outputFile string, cfg internal.Config) error) error {
	workChan := make(chan string)
	errorChan := make(chan error, len(files))

	wg := &sync.WaitGroup{}
	wg.Add(parallelism)

	for i := 0; i < parallelism; i++ {
		go func() {
			defer wg.Done()
			for file := range workChan {
				if err := fn(file, file, Cfg); err != nil {
					errorChan <- err
				}
			}
		}()
	}
	sliceLen := len(files)
	for i := 0; i < sliceLen; i++ {
		workChan <- files[i]
	}

	close(workChan)
	wg.Wait()
	close(errorChan)

	var processingErrors []error
	for err := range errorChan {
		processingErrors = append(processingErrors, err)
	}
	return errors.Join(processingErrors...)
}
