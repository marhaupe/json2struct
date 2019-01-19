package cmd

import (
	"fmt"
	"os"

	"github.com/marhaupe/json2struct/internal"
	"github.com/marhaupe/json2struct/internal/editor"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var rootCmd = &cobra.Command{
	Use:     "json2struct <jsonString>",
	Short:   "Generate a struct from a JSON string argument",
	Example: "json2struct \"$(curl \"https://reqres.in/api/users?page=2\")\"",
	Args:    cobra.ExactArgs(1),

	Run: rootFunc,
}

func rootFunc(cmd *cobra.Command, args []string) {
	jsonstr := args[0]
	gen, err := internal.Generate(jsonstr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(gen)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create your JSON in the editor",
	Long:  "Create your JSON in the editor. Make sure to save the file via :wq or similar",
	Args:  cobra.ExactArgs(0),

	Run: createFunc,
}

func createFunc(cmd *cobra.Command, args []string) {
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
	fmt.Println(gen)
}
