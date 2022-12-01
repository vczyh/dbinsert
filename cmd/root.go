package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbinsert",
	Short: "A quick insert tool, support mysql, postgresql.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	definitionFile string
	tableSize      int
)

func init() {
	rootCmd.PersistentFlags().StringVar(&definitionFile, "definition", "", "definition file path")
}