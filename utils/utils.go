package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/fsouza/go-dockerclient"
)

var docker_client *docker.Client

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return (err == nil)
}

func GetCommandExitCode(err error) int {
	var waitStatus syscall.WaitStatus
	if exitError, ok := err.(*exec.ExitError); ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
		return waitStatus.ExitStatus()
	}
	return 0
}

func GetFileContent(path string) string {
	var content string
	content_bytes, err := ioutil.ReadFile(path)
	content = string(content_bytes)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return content
}

func WriteFileContent(filename string, content string) {
	content_bytes := []byte(content)
	err := ioutil.WriteFile(filename, content_bytes, 0644)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DockerImageExists(name string) bool {
	if docker_client == nil {
		docker_client, _ = docker.NewClientFromEnv()
	}
	if _, err := docker_client.InspectImage(name); err == docker.ErrNoSuchImage {
		return false
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return true
}

func DockerRemoveImage(name string) bool {
	if docker_client == nil {
		docker_client, _ = docker.NewClientFromEnv()
	}
	err := docker_client.RemoveImage(name)
	if err == docker.ErrNoSuchImage {
		return false
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return true
}

func RunCmd(command string) (int, string) {
	exitcode := 0
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		exitcode = GetCommandExitCode(err)
	}
	return exitcode, string(out)
}
