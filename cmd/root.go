package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "provide a JSON as first arg to generate a struct",
	Short: "Generate a struct from a JSON document",
	Long: `json2struct generates a struct from a JSON document. 
	Visit https://github.com/marhaupe/json2struct for documentation and for issues`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gen := Generate(args[0])
		fmt.Println(gen)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
