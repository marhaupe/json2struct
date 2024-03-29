package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/marhaupe/json2struct/pkg/editor"
	"github.com/marhaupe/json2struct/pkg/generator"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

var (
	inputString        string
	inputFile          string
	version            string
	shouldBenchmark    bool
	shouldUseClipboard bool

	rootCmd = &cobra.Command{
		Use:     "json2struct",
		Short:   "json2struct generates Go type definitions for a JSON",
		Version: version,
		Args:    cobra.ExactArgs(0),
		Run:     rootFunc,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&inputString, "string", "s", "", "JSON string")
	rootCmd.Flags().StringVarP(&inputFile, "file", "f", "", "path to JSON file")
	rootCmd.Flags().BoolVarP(&shouldBenchmark, "benchmark", "b", false, "measure execution time")
	rootCmd.Flags().BoolVarP(&shouldUseClipboard, "clipboard", "c", false, "read from and write types to clipboard")
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
	case shouldUseClipboard:
		err := clipboard.Init()
		if err != nil {
			os.Exit(2)
		}
		userInput = string(clipboard.Read(clipboard.FmtText))
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

	output, err := generator.GenerateOutputFromString(userInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	if shouldUseClipboard {
		err := clipboard.Init()
		if err != nil {
			os.Exit(4)
		}
		clipboard.Write(clipboard.FmtText, []byte(output))
		fmt.Println("saved output to clipboard")
	} else {
		fmt.Println(output)
	}
}

func readFromFile() string {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
