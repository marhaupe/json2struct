package cmd

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"

	"github.com/marhaupe/json2struct/internal/editor"
	"github.com/marhaupe/json2struct/internal/generate"
	"github.com/spf13/cobra"
)

// JSONString may optionally contain the JSON provided by the user
var JSONString string

// JSONFile may optionally contain the path to a JSON file provided by the user
var JSONFile string

func init() {
	rootCmd.Flags().StringVarP(&JSONString, "string", "s", "", "JSON string")
	rootCmd.Flags().StringVarP(&JSONFile, "file", "f", "", "Path to JSON file")
}

var rootCmd = &cobra.Command{
	Use:   "json2struct",
	Short: "These are all available commands to help you parse JSONs to Go structs",
	Args:  cobra.ExactArgs(0),
	Run:   rootFunc,
}

func rootFunc(cmd *cobra.Command, args []string) {
	var res string
	if JSONFile != "" {
		res = generateFromFile()
	} else if JSONString != "" {
		res = generateFromString()
	} else {
		res = generateFromEditor()
	}
	format, err := format.Source([]byte(res))
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	fmt.Println(string(format))
}

func generateFromFile() string {
	data, err := ioutil.ReadFile(JSONFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	gen, err := generate.Generate(string(data))
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	return gen
}

func generateFromString() string {
	gen, err := generate.Generate(JSONString)
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
	gen, err := generate.Generate(jsonstr)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	return gen
}
