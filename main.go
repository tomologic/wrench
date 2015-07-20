package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "wrench"}

	main_config(rootCmd)
	main_build(rootCmd)

	rootCmd.Execute()
}
