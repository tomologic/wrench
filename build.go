package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func main_build(rootCmd *cobra.Command) {
	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "Build docker image",
		Long:  `will build docker image for project`,
		Run: func(cmd *cobra.Command, args []string) {
			build()
		},
	}

	rootCmd.AddCommand(cmdBuild)
}

func build() {
	if fileExists("./Dockerfile.builder") {
		fmt.Printf("INFO: %s\n", "Builder build mode")
		buildBuilder()
	} else if fileExists("./Dockerfile") {
		fmt.Printf("INFO: %s\n", "Simple build mode")
		buildSimple()
	} else {
		fmt.Printf("ERROR: %s\n", "No Dockerfile found.")
		os.Exit(1)
	}
}

func buildBuilder() {
	image_name := fmt.Sprintf("%s/%s:%s",
		config.Project.Organization,
		config.Project.Name,
		config.Project.Version)

	builder_image_name := fmt.Sprintf("%s/builder-%s:%s",
		config.Project.Organization,
		config.Project.Name,
		config.Project.Version)

	fmt.Printf("INFO: %s %s\n\n",
		"Found Dockerfile.builder, building image builder",
		builder_image_name)

	// Build builder image
	cmd_string := fmt.Sprintf("docker build -f Dockerfile.builder -t '%s' .", builder_image_name)
	cmd := exec.Command("sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("\nINFO: %s %s\n\n",
		"Building image with builder",
		image_name)

	// Build image
	cmd_string = fmt.Sprintf("docker run --rm '%s' | docker build -t '%s' -", builder_image_name, image_name)
	cmd = exec.Command("sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Cleanup builder image
	cmd_string = fmt.Sprintf("docker rmi '%s'", builder_image_name)
	cmd = exec.Command("sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildSimple() {
	image_name := fmt.Sprintf("%s/%s:%s",
		config.Project.Organization,
		config.Project.Name,
		config.Project.Version)

	fmt.Printf("INFO: %s %s\n",
		"Found Dockerfile, building image",
		image_name)

	cmd_string := fmt.Sprintf("docker build -t '%s' .", image_name)
	cmd := exec.Command("sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
