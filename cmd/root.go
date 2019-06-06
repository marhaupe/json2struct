package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/marhaupe/json2struct/internal/editor"
	"github.com/spf13/cobra"
)

var (
	inputString string
	inputFile   string
	version     string

	rootCmd = &cobra.Command{
		Use:     "json2struct",
		Short:   "Parse a JSON into a generated Go struct",
		Version: version,
		Args:    cobra.ExactArgs(0),
		Run:     rootFunc,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&inputString, "string", "s", "", "JSON string")
	rootCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Path to JSON file")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func rootFunc(cmd *cobra.Command, args []string) {
	var json string
	switch {
	case inputFile != "":
		json = readFromFile()
	case inputString != "":
		json = inputString
	default:
		json = readFromEditor()
	}

	res := generate(json)

	fmt.Println(res)
}

func readFromFile() string {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	return string(data)
}

func readFromEditor() string {
	return awaitValidInput()
}

func awaitValidInput() string {
	edit := editor.New()
	defer edit.Delete()
	edit.Display()

	var jsonstr string
	jsonstr, _ = edit.Read()

	isValid := json.Valid([]byte(jsonstr))
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
		isValid := json.Valid([]byte(jsonstr))
		if isValid {
			return jsonstr
		}
	}
}

func generate(json string) string {
	// node := parse.ParseFromString(json)

	// ast := generator.GenerateAST(node)

	// fileset := token.NewFileSet()
	// var buf bytes.Buffer
	// printer.Fprint(&buf, fileset, ast)

	// return fmt.Sprintf("%s\n", buf.String())
	return ""
}
