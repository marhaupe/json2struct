package cmd

import (
	"fmt"
	"os"

	"github.com/marhaupe/json2struct/internal"
	"github.com/marhaupe/json2struct/internal/editor"
	"github.com/spf13/cobra"
)

// JSONString may optionally contain the JSON provided by the user
var JSONString string

func init() {
	rootCmd.Flags().StringVarP(&JSONString, "string", "s", "", "JSON file as string")
}

var rootCmd = &cobra.Command{
	Use:   "json2struct",
	Short: "These are all available commands to help you parse JSONs to Go structs",
	Args:  cobra.ExactArgs(0),
	Run:   rootFunc,
}

func rootFunc(cmd *cobra.Command, args []string) {
	if JSONString != "" {
		fmt.Println(generateFromString())
	} else {
		fmt.Println(generateFromEditor())
	}
}

func generateFromString() string {
	gen, err := internal.Generate(JSONString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return gen
}

func generateFromEditor() string {
	editor := editor.New()
	editor.Display()
	jsonstr, err := editor.Consume()
	if err != nil {
		fmt.Println("Error while reading from VIM", err)
		os.Exit(2)
	}
	gen, err := internal.Generate(jsonstr)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	return gen
}
