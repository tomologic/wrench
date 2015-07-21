package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "unknown"

func AddToWrench(cmdRoot *cobra.Command) {
	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Version of wrench",
		Long:  `Version of wrench`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	cmdRoot.AddCommand(cmdVersion)
}
