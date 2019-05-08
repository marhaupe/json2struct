package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/marhaupe/json2struct/internal/editor"
	"github.com/marhaupe/json2struct/internal/generate"
	"github.com/marhaupe/json2struct/internal/lex"
	"github.com/spf13/cobra"
)

var version string

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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Args:  cobra.ExactArgs(0),
	Run:   cmdFunc,
}

func init() {
	rootCmd.Flags().StringVarP(&InputString, "string", "s", "", "JSON string")
	rootCmd.Flags().StringVarP(&InputFile, "file", "f", "", "Path to JSON file")
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func cmdFunc(cmd *cobra.Command, args []string) {
	fmt.Println(version)
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
	jsonstr := awaitValidInput()
	gen, err := generate.GenerateWithFormatting(jsonstr)
	if err != nil {
		os.Exit(2)
	}
	return gen
}

func awaitValidInput() string {
	edit := editor.New()
	defer edit.Delete()
	edit.Display()

	var jsonstr string
	jsonstr, _ = edit.Read()

	isValid := lex.ValidateJSON(jsonstr)
	if isValid {
		return jsonstr
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You supplied an invalid json. Do you want to fix it (y/n)?  ")

		input, _ := reader.ReadString('\n')
		userWantsFix := string(input[0]) == "y"
		if !userWantsFix {
			return ""
		}

		edit.Display()
		jsonstr, _ = edit.Read()
		isValid := lex.ValidateJSON(jsonstr)
		if isValid {
			return jsonstr
		}
	}
}
