package cmd

import (
	"fmt"
	"runtime"
	"sync"

	internal "github.com/ffais/yaml-sort/internal"
	"github.com/spf13/cobra"
)

var OutputFile string

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Yaml-Sort sorts content of YAML files alphabetically.",
	Long:  `Yaml-Sort sorts content of YAML files alphabetically preserving comments, anchor and with support for custom order.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		searchDir, _ := cmd.Flags().GetString("search-dir")
		if searchDir == "" {
			cmd.MarkFlagRequired("output-file")
			cmd.MarkPersistentFlagRequired("input-file")
		}
	},
	Run: sort,
}

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.Flags().StringVarP(&OutputFile, "output-file", "o", "", "The YAML file to output sorted content to.")
}

func sort(cmd *cobra.Command, args []string) {
	if Cfg.SearchDir != "" {
		parallelism := runtime.NumCPU() * 2
		yamls, _ := internal.FindYamlFile(Cfg.SearchDir, InputFile)
		parallelProcessing(yamls, parallelism, sortYamlFile)
	} else {
		sortYamlFile(InputFile, OutputFile, Cfg)
	}
}

func sortYamlFile(inputFile string, outputFile string, cfg internal.Config) {
	fmt.Println("Sorting yaml file", inputFile)
	documents := internal.ParseYaml(inputFile)
	for _, document := range documents {
		internal.SortYamlNodes(document, cfg)
		if Cfg.SpaceTopKey {
			internal.AddEmptyLinesBeforeTopLevelKeys(document)
		}
	}
	internal.WriteToFile(outputFile, documents, cfg)
}

func parallelProcessing(files []string, parallelism int, fn func(inputFile string, outputFile string, cfg internal.Config)) {
	workChan := make(chan string)

	wg := &sync.WaitGroup{}
	wg.Add(parallelism)

	for i := 0; i < parallelism; i++ {
		go func() {
			defer wg.Done()
			for file := range workChan {
				fn(file, file, Cfg)
			}
		}()
	}
	sliceLen := len(files)
	for i := 0; i < sliceLen; i++ {
		workChan <- files[i]
	}

	close(workChan)
	wg.Wait()
}
