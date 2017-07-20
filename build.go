package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomologic/wrench/config"
	"github.com/tomologic/wrench/utils"
)

var flag_rebuild bool

func AddBuildToWrench(rootCmd *cobra.Command) {
	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "Build docker image",
		Long:  `will build docker image for project`,
		Run: func(cmd *cobra.Command, args []string) {
			build()
		},
	}

	cmdBuild.Flags().BoolVarP(&flag_rebuild, "rebuild", "r", false, "Force rebuild of image")
	rootCmd.AddCommand(cmdBuild)
}

func build() {
	image_name := config.GetProjectImage()

	if !flag_rebuild && utils.DockerImageExists(image_name) {
		fmt.Printf("INFO: Docker image %s already exists\n", image_name)

		// Build test image if missing
		if !utils.DockerImageExists(fmt.Sprintf("%s-test", image_name)) {
			buildTest()
		}

		os.Exit(0)
	}

	if utils.FileExists("./Dockerfile.builder") {
		buildBuilder()
		buildTest()
	} else if utils.FileExists("./Dockerfile") {
		buildSimple()
		buildTest()
	} else {
		fmt.Printf("ERROR: %s\n", "No Dockerfile found.")
		os.Exit(1)
	}
}

func buildBuilder() {
	image_name := config.GetProjectImage()

	builder_image_name := fmt.Sprintf("%s-builder", image_name)

	fmt.Printf("INFO: %s %s\n\n",
		"Found Dockerfile.builder, building image builder",
		builder_image_name)

	// Build builder image
	cmd_string := fmt.Sprintf("docker build -f Dockerfile.builder -t '%s' .", builder_image_name)
	cmd := exec.Command("sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	version := strings.TrimLeft(config.GetProjectVersion(), "v")
	fmt.Printf("INFO: Adding env variable VERSION=%s\n\n", version)
	utils.DockerImageAddEnv(builder_image_name, "VERSION", version)

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

	fmt.Printf("INFO: Adding env variable VERSION=%s\n\n", version)
	utils.DockerImageAddEnv(image_name, "VERSION", version)
}

func buildSimple() {
	image_name := config.GetProjectImage()

	fmt.Printf("INFO: %s %s\n\n",
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

	version := strings.TrimLeft(config.GetProjectVersion(), "v")
	fmt.Printf("INFO: Adding env variable VERSION=%s\n\n", version)
	utils.DockerImageAddEnv(image_name, "VERSION", version)
}

func buildTest() {
	image_name := config.GetProjectImage()

	test_image_name := fmt.Sprintf("%s-test", image_name)

	if !utils.FileExists("./Dockerfile.test") {
		return
	}

	fmt.Printf("INFO: %s %s\n\n",
		"Found Dockerfile.test, building test image",
		test_image_name)

	dockerfile := utils.GetFileContent("./Dockerfile.test")

	dockerfile_lines := strings.Split(dockerfile, "\n")

	if !strings.HasPrefix(dockerfile_lines[0], "FROM") {
		fmt.Println("ERROR: Missing FROM on first line in Dockerfile.test")
		os.Exit(1)
	}

	// if FROM string subfix with builder then base on builder image
	if strings.HasSuffix(dockerfile_lines[0], "builder") {
		dockerfile_lines[0] = fmt.Sprintf("FROM %s-builder", image_name)
	} else {
		dockerfile_lines[0] = fmt.Sprintf("FROM %s", image_name)
	}

	temp_dockerfile_content := strings.Join(dockerfile_lines, "\n")

	// Tempdir for building test image
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	tempdir, err := ioutil.TempDir(dir, ".wrench_build_")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempdir)

	temp_dockerfile := fmt.Sprintf("%s/Dockerfile.test", tempdir)
	utils.WriteFileContent(temp_dockerfile, temp_dockerfile_content)

	cmd_string := fmt.Sprintf("docker build -f %s -t '%s' .", temp_dockerfile, test_image_name)
	cmd := exec.Command("sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		os.RemoveAll(tempdir)
		fmt.Println(err)
		os.Exit(1)
	}
}
