package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomologic/wrench/build"
	"github.com/tomologic/wrench/bump"
	"github.com/tomologic/wrench/config"
	"github.com/tomologic/wrench/push"
	"github.com/tomologic/wrench/run"
)

var VERSION = "0.0.0"

func main() {
	var rootCmd = &cobra.Command{Use: "wrench"}

	build.AddToWrench(rootCmd)
	bump.AddToWrench(rootCmd)
	push.AddToWrench(rootCmd)
	config.AddToWrench(rootCmd)
	run.AddToWrench(rootCmd)

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Version of wrench",
		Long:  `Version of wrench`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(strings.Trim(VERSION, "'"))
		},
	}
	rootCmd.AddCommand(cmdVersion)

	rootCmd.Execute()
}
