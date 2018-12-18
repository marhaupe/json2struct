package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "json2struct",
	Short: "Generate a struct from a JSON document",
	Long: "json2struct generates a struct from a JSON document.\n" +
		"Visit https://github.com/marhaupe/json2struct for further documentation.\n" +
		"Feel free to open an issue if you encounter any bugs!",
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			jsonstr, err := internal.VimToString("json2struct.temp")
			if err != nil {
				panic(err)
			}
			gen := Generate(jsonstr)
			fmt.Println(gen)
		case 1:
			jsonstr := args[0]
			gen := Generate(jsonstr)
			fmt.Println(gen)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
