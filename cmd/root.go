package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/marhaupe/json2struct/pkg/editor"
	"github.com/marhaupe/json2struct/pkg/generator"
	"github.com/marhaupe/json2struct/pkg/parse"
	"github.com/spf13/cobra"
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
		Run: func(cmd *cobra.Command, args []string) {
			Run(inputString, inputFile, shouldBenchmark, shouldUseClipboard)
		},
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

func Run(inputString, inputFile string, shouldBenchmark, shouldUseClipboard bool) {
	if shouldBenchmark {
		defer benchmark()()
	}

	var userInputNode parse.Node
	var err error

	switch {
	case shouldUseClipboard:
		var userInput string
		userInput, err = clipboard.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		userInputNode, err = parse.ParseFromString(userInput)
	case inputFile != "":
		userInput := readFromFile()
		userInputNode, err = parse.ParseFromString(userInput)
	case inputString != "":
		userInput := inputString
		userInputNode, err = parse.ParseFromString(userInput)
	default:
		userInputNode, err = readFromEditor()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if userInputNode == nil {
		return
	}

	output, err := generator.GenerateOutputFromAST(userInputNode)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	fmt.Println(output)

	if shouldUseClipboard {
		err = clipboard.WriteAll(output)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
		fmt.Println("\nSaved output to clipboard")
	}
}

func readFromFile() string {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(data)
}

func readFromEditor() (parse.Node, error) {
	edit := editor.New()
	defer edit.Delete()
	edit.Display()

	userInput, _ := edit.Read()

	userInputNode, err := parse.ParseFromString(userInput)
	if err == nil {
		return userInputNode, nil
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You supplied invalid JSON. Continue editing? (Y/n) ")
		userAnswer, _ := reader.ReadString('\n')
		userAnswer = strings.TrimSpace(userAnswer)
		userWantsFix := len(userAnswer) == 0 || userAnswer[0] == 'y'
		if !userWantsFix {
			return nil, nil
		}
		fmt.Print("\033[1A\033[2K")
		edit.Display()
		userInput, _ = edit.Read()
		isValid := json.Valid([]byte(userInput))
		if isValid {
			return userInputNode, nil
		}
	}
}

func benchmark() func() {
	start := time.Now()
	return func() {
		fmt.Printf("generating took %v\n", time.Since(start))
	}
}
