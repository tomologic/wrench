package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

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

func GetFileContent(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return lines
}

func WriteFileContent(filename string, lines []string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		file.Close()
		os.Exit(1)
	}

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	if err = w.Flush(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
