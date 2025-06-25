package cmd

import (
	"fmt"
	"os"

	internal "github.com/ffais/yaml-sort/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, InputFile, SearchDir string
var customSort []string
var reverse, spaceTopKey bool
var indent int
var Cfg internal.Config

var rootCmd = &cobra.Command{
	Use:              "yaml-sort",
	Short:            "yaml-sort sort yaml file presernvig comments and with custom order",
	Long:             `yaml-sort sort yaml file presernvig comments and with custom order`,
	PersistentPreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Yaml-Sort! Use --help to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&InputFile, "input-file", "i", "", "input file")
	rootCmd.PersistentFlags().StringSliceVarP(&customSort, "custom-sort", "c", []string{}, "Custom-Order")
	rootCmd.PersistentFlags().BoolVarP(&reverse, "reverse", "r", false, "Reverse the order")
	rootCmd.PersistentFlags().BoolVarP(&spaceTopKey, "space-top-key", "s", true, "Add an empty line beetween top level keys")
	rootCmd.PersistentFlags().IntVarP(&indent, "indent", "t", 2, "Reverse the order")
	rootCmd.PersistentFlags().StringVarP(&SearchDir, "search-dir", "d", "", "search dir")
	rootCmd.MarkPersistentFlagRequired("input-file")
	viper.BindPFlag("custom-sort", rootCmd.PersistentFlags().Lookup("custom-sort"))
	viper.BindPFlag("reverse", rootCmd.PersistentFlags().Lookup("reverse"))
	viper.BindPFlag("space-top-key", rootCmd.PersistentFlags().Lookup("space-top-key"))
	viper.BindPFlag("indent", rootCmd.PersistentFlags().Lookup("indent"))
	viper.BindPFlag("search-dir", rootCmd.PersistentFlags().Lookup("search-dir"))
}

func initConfig(cmd *cobra.Command, args []string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				println("no file")
			} else {
				// Config file was found but another error was produced
				panic(fmt.Errorf("fatal error config file: %s", err))
			}
		}
	}
	err := viper.Unmarshal(&Cfg)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
}
