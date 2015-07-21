package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomologic/wrench/utils"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Organization *string
	Name         *string
	Version      *string
}
type Config struct {
	Project Project
}

var config = &Config{}

func AddToWrench(cmdRoot *cobra.Command) {
	readWrenchFile()

	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Configuration for wrench",
		Long:  `configuration picked up by wrench and used in commands`,
		Run: func(cmd *cobra.Command, args []string) {
			generateAllConfig()
			d, err := yaml.Marshal(&config)
			if err != nil {
				panic(err)
			}
			fmt.Printf(string(d))
		},
	}

	cmdRoot.AddCommand(cmdConfig)
}

func readWrenchFile() {
	if !utils.FileExists("./wrench.yml") {
		return
	}

	file, err := ioutil.ReadFile("./wrench.yml")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateAllConfig() {
	GetProjectOrganization()
	GetProjectName()
	GetProjectVersion()
}

func GetProjectOrganization() string {
	if config.Project.Organization == nil {
		config.Project.Organization = detectProjectOrganization()
	}
	return *config.Project.Organization
}

func GetProjectName() string {
	if config.Project.Name == nil {
		config.Project.Name = detectProjectName()
	}
	return *config.Project.Name
}

func GetProjectVersion() string {
	if config.Project.Version == nil {
		config.Project.Version = detectProjectVersion()
	}
	return *config.Project.Version
}

func detectProjectOrganization() *string {
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
	org = strings.TrimSpace(org)
	return &org
}

func detectProjectName() *string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	project := string(filepath.Base(dir))
	return &project
}

func detectProjectVersion() *string {
	out, err := exec.Command("sh", "-c", "git describe").Output()
	if err != nil {
		panic(err)
	}
	version := strings.TrimSpace(string(out))
	return &version
}
