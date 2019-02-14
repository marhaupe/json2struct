// Package cmd contains fields to configure the cli tool shipped with json2struct
package cmd

import "fmt"

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
