package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version        = "dev"
	commitHash     = "n/a"
	buildTimestamp = "n/a"
)

func BuildVersion() string {
	return fmt.Sprintf("version.BuildInfo{Version:\"%s\", GitCommit:\"%s\", BuildTimestamp:\"%s\"}", version, commitHash, buildTimestamp)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number, commit hash and build timestamp of Yaml-Sort",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(BuildVersion())
	},
}
