package main

import (
	"github.com/spf13/cobra"
	"github.com/tomologic/wrench/build"
	"github.com/tomologic/wrench/config"
	"github.com/tomologic/wrench/run"
	"github.com/tomologic/wrench/version"
)

func main() {
	var rootCmd = &cobra.Command{Use: "wrench"}

	build.AddToWrench(rootCmd)
	config.AddToWrench(rootCmd)
	run.AddToWrench(rootCmd)
	version.AddToWrench(rootCmd)

	rootCmd.Execute()
}
