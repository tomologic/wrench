package utils

import (
	"archive/tar"
	"bytes"
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

func DockerImageAddEnv(image, env, value string) error {
	var files = make([]Tarfile, 2)
	files[0] = Tarfile{
		"Dockerfile",
		fmt.Sprintf("FROM %s\nENV %s %s\n", image, env, value),
	}
	tarfile, err := CreateTar(files)
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "build", "-t", image, "-")

	// Open the tar archive for reading.
	cmd.Stdin = bytes.NewReader(tarfile.Bytes())

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func RunCmd(command string) (int, string) {
	exitcode := 0
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		exitcode = GetCommandExitCode(err)
	}
	return exitcode, string(out)
}

type Tarfile struct {
	Name, Content string
}

func CreateTar(files []Tarfile) (*bytes.Buffer, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new tar archive.
	tw := tar.NewWriter(buf)

	// Add some files to the archive.
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(file.Content)); err != nil {
			return nil, err
		}
	}
	// Make sure to check the error on Close.
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
