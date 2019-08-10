package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/marhaupe/json2struct/internal/editor"
	"github.com/marhaupe/json2struct/internal/generator"
	"github.com/marhaupe/json2struct/internal/parse"
	"github.com/spf13/cobra"
)

var (
	inputString     string
	inputFile       string
	version         string
	shouldBenchmark bool

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
	rootCmd.Flags().BoolVarP(&shouldBenchmark, "benchmark", "b", false, "Measure execution time")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func rootFunc(cmd *cobra.Command, args []string) {
	var userInput string
	switch {
	case inputFile != "":
		userInput = readFromFile()
	case inputString != "":
		userInput = inputString
	default:
		userInput = readFromEditor()
	}

	if shouldBenchmark {
		defer benchmark()()
	}

	generatedCode, err := generate(userInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(generatedCode)
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
	edit := editor.New()
	defer edit.Delete()
	edit.Display()

	var userInput string
	userInput, _ = edit.Read()

	isValid := json.Valid([]byte(userInput))
	if isValid {
		return userInput
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You supplied an invalid JSON. Do you want to fix it (y/n)?\t")
		userAnswer, _ := reader.ReadString('\n')
		userWantsFix := string(userAnswer[0]) == "y"
		if !userWantsFix {
			return ""
		}
		edit.Display()
		userInput, _ = edit.Read()
		isValid := json.Valid([]byte(userInput))
		if isValid {
			return userInput
		}
	}
}

func benchmark() func() {
	start := time.Now()
	return func() {
		fmt.Printf("generating took %v\n", time.Since(start))
	}
}

func generate(json string) (string, error) {
	node, err := parse.ParseFromString("json2struct", json)
	if err != nil {
		return "", err
	}

	file, err := generator.GenerateFile(node)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = file.Render(buf)
	return buf.String(), err
}
