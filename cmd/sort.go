package cmd

import (
	"fmt"

	internal "github.com/ffais/yaml-sort/internal"
	"github.com/spf13/cobra"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var OutputFile string

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.Flags().StringVarP(&OutputFile, "output-file", "o", "", "output file")
	sortCmd.MarkFlagRequired("output-file")
}

func sort(cmd *cobra.Command, args []string) {
	var node yaml.Node
	fmt.Println("Sorting yaml file", InputFile, OutputFile)
	internal.ParseYaml(InputFile, &node)
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		rootNode := node.Content[0]
		internal.SortYamlNodes(rootNode, Cfg)
		if Cfg.SpaceTopKey {
			internal.AddEmptyLinesBeforeTopLevelKeys(rootNode)
		}
	}

	internal.WriteToFile(OutputFile, &node, Cfg)
}

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort yaml file",
	Long:  "Sort yaml file",
	Run:   sort,
}
