package cmd

import "fmt"

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
