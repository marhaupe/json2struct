package cmd

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

var (
	inputString string
	inputFile   string
)

var rootCmd = &cobra.Command{
	Use:   "json2struct",
	Short: "These are all available commands to help you parse JSONs to Go structs",
	Args:  cobra.ExactArgs(0),
	Run:   rootFunc,
}

func init() {
	rootCmd.Flags().StringVarP(&inputString, "string", "s", "", "JSON string")
	rootCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Path to JSON file")
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func rootFunc(cmd *cobra.Command, args []string) {
	var res string

	switch {
	case inputFile != "":
		res = generateFromFile()
	case inputString != "":
		res = generateFromString()
	default:
		res = generateFromEditor()
	}

	fmt.Println(res)
}

func generateFromFile() string {
	data, err := ioutil.ReadFile(inputFile)
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
	gen, err := generate.GenerateWithFormatting(inputString)
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
