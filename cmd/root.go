package cmd

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbinsert",
	Short: "A quick insert tool, support mysql, postgresql, redis.",
}

func Execute(aVersion string) {
	version = aVersion

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	definitionFile string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&definitionFile, "schema", "", "definition file path")
}
