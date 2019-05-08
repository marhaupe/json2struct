package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Args:  cobra.ExactArgs(0),
	Run:   versionFunc,
}

func versionFunc(cmd *cobra.Command, args []string) {
	fmt.Println(version)
}
