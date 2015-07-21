package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Organization string
	Name         string
	Version      string
}
type Config struct {
	Project *Project
}

var config *Config

func main_config(cmdRoot *cobra.Command) {
	GenerateConfig()

	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Configuration for wrench",
		Long:  `configuration picked up by wrench and used in commands`,
		Run: func(cmd *cobra.Command, args []string) {
			d, err := yaml.Marshal(&config)
			if err != nil {
				panic(err)
			}
			fmt.Printf(string(d))
		},
	}

	cmdRoot.AddCommand(cmdConfig)
}

func GenerateConfig() {
	var project = &Project{
		Organization: detectProjectOrganization(),
		Name:         detectProjectName(),
		Version:      detectProjectVersion(),
	}

	config = &Config{
		Project: project}
}

func detectProjectOrganization() string {
	out, err := exec.Command("sh", "-c", "hostname -f").Output()
	if err != nil {
		panic(err)
	}
	parts := strings.Split(string(out), ".")

	var org string
	if len(parts) <= 2 {
		// handle user.local
		org = parts[len(parts)-1]
	} else {
		// handle user.organization.com
		org = parts[len(parts)-2]
	}
	return strings.TrimSpace(org)
}

func detectProjectName() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return string(filepath.Base(dir))
}

func detectProjectVersion() string {
	out, err := exec.Command("sh", "-c", "git describe").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}
