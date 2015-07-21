package main

import (
	"github.com/spf13/cobra"
	"github.com/tomologic/wrench/build"
	"github.com/tomologic/wrench/config"
)

func main() {
	var rootCmd = &cobra.Command{Use: "wrench"}

	config.AddToWrench(rootCmd)
	build.AddToWrench(rootCmd)

	rootCmd.Execute()
}
