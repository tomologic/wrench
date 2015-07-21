package utils

import (
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
