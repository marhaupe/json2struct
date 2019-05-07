package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/marhaupe/json2struct/internal/editor"
	"github.com/marhaupe/json2struct/internal/generate"
	"github.com/spf13/cobra"
)

// InputString may optionally contain the JSON provided by the user
var InputString string

// InputFile may optionally contain the path to a JSON file provided by the user
var InputFile string

var rootCmd = &cobra.Command{
	Use:   "json2struct",
	Short: "These are all available commands to help you parse JSONs to Go structs",
	Args:  cobra.ExactArgs(0),
	Run:   rootFunc,
}

func init() {
	rootCmd.Flags().StringVarP(&InputString, "string", "s", "", "JSON string")
	rootCmd.Flags().StringVarP(&InputFile, "file", "f", "", "Path to JSON file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func rootFunc(cmd *cobra.Command, args []string) {
	var res string

	switch {
	case InputFile != "":
		res = generateFromFile()
	case InputString != "":
		res = generateFromString()
	default:
		res = generateFromEditor()
	}

	fmt.Println(res)
}

func generateFromFile() string {
	data, err := ioutil.ReadFile(InputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	gen, err := generate.GenerateWithFormatting(string(data))
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	return gen
}

func generateFromString() string {
	gen, err := generate.GenerateWithFormatting(InputString)
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
	gen, err := generate.GenerateWithFormatting(jsonstr)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	return gen
}
