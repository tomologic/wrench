package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/tomologic/wrench/utils"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Organization string `yaml:"Organization"`
	Name         string `yaml:"Name"`
	Version      string `yaml:"Version"`
	Image        string `yaml:"Image"`
}
type Run struct {
	Cmd string   `yaml:"Cmd"`
	Env []string `yaml:"Env,omitempty"`
}
type Config struct {
	Project Project        `yaml:"Project"`
	Run     map[string]Run `yaml:"Run,omitempty"`
}

var config = &Config{}
var flag_format string

func AddToWrench(cmdRoot *cobra.Command) {
	readWrenchFile()

	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Configuration for wrench",
		Long:  `configuration picked up by wrench and used in commands`,
		Run: func(cmd *cobra.Command, args []string) {
			commandConfig()
		},
	}

	cmdConfig.Flags().StringVar(&flag_format, "format", "", "Return specific value from config")

	cmdRoot.AddCommand(cmdConfig)
}

func commandConfig() {
	generateAllConfig()

	if flag_format == "" {
		d, err := yaml.Marshal(&config)
		if err != nil {
			panic(err)
		}
		fmt.Printf(string(d))
	} else {
		tmpl, err := template.New("format").Parse(flag_format)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		err = tmpl.Execute(os.Stdout, &config)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	}
}

func readWrenchFile() {
	if !utils.FileExists("./wrench.yml") {
		return
	}

	// Get wrench file content
	file, err := ioutil.ReadFile("./wrench.yml")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	// Create structure accessible from wrench file
	Environ := make(map[string]string)
	type TemplateContext struct {
		Environ *map[string]string
	}
	tmpl_context := TemplateContext{
		Environ: &Environ,
	}

	// Get all environment variables
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		Environ[splits[0]] = strings.Join(splits[1:], "=")
	}

	// Create template from wrench file
	var rendered_config bytes.Buffer
	tmpl, err := template.New("config").Parse(string(file))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Render template with tmpl_context
	err = tmpl.Execute(&rendered_config, tmpl_context)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	type UnmarshalConfig struct {
		Project Project       `yaml:"Project"`
		Run     yaml.MapSlice `yaml:"Run,omitempty"`
	}
	uconfig := UnmarshalConfig{}

	// Load the expected yaml file structure from rendered config
	err = yaml.Unmarshal(rendered_config.Bytes(), &uconfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get Project from unmarshalled config
	config.Project = uconfig.Project
	config.Run = make(map[string]Run)

	// Handle errors parsing dynamic structure of Run
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERROR: Unexpected yaml structure for Run\n")
			os.Exit(1)
		}
	}()

	// Handle dynamic structure of Run map
	for _, item := range uconfig.Run {
		name, ok := item.Key.(string)
		if ok != true {
			fmt.Println("ERROR: Unexpected yaml structure for Run")
			os.Exit(1)
		}

		run := Run{}

		// No more values provided if value is a string
		run.Cmd, ok = item.Value.(string)
		if ok {
			config.Run[name] = run
			continue
		}

		r := item.Value.(yaml.MapSlice)
		for k := range r {
			if r[k].Key.(string) == "Cmd" {
				run.Cmd = r[k].Value.(string)
			} else if r[k].Key.(string) == "Env" {
				for _, a := range r[k].Value.([]interface{}) {
					run.Env = append(run.Env, a.(string))
				}
			} else {
				panic(fmt.Sprintf("Unknown key %s\n", r[k].Key.(string)))
			}
		}

		if run.Cmd == "" {
			fmt.Printf("ERROR: Unexpected yaml structure for Run.%s\n", name)
			os.Exit(1)
		}

		config.Run[name] = run
	}
}

func generateAllConfig() {
	GetProjectOrganization()
	GetProjectName()
	GetProjectVersion()
	GetProjectImage()
}

func GetProjectOrganization() string {
	if config.Project.Organization == "" {
		config.Project.Organization = detectProjectOrganization()
	}
	return config.Project.Organization
}

func GetProjectName() string {
	if config.Project.Name == "" {
		config.Project.Name = detectProjectName()
	}
	return config.Project.Name
}

func GetProjectVersion() string {
	if config.Project.Version == "" {
		config.Project.Version = detectProjectVersion()
	}
	return config.Project.Version
}

func GetProjectImage() string {
	if config.Project.Image == "" {
		config.Project.Image = fmt.Sprintf("%s/%s:%s",
			GetProjectOrganization(),
			GetProjectName(),
			GetProjectVersion())
	}
	return config.Project.Image
}

func GetRun(name string) (Run, bool) {
	val, ok := config.Run[name]
	return val, ok
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
	org = strings.TrimSpace(org)
	return org
}

func detectProjectName() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	project := string(filepath.Base(dir))
	return project
}

func detectProjectVersion() string {
	// make sure git is installed and we are inside a git repo
	cmd := exec.Command("sh", "-c", "git rev-parse --short HEAD")
	out, err := cmd.Output()
	if err != nil {
		exitcode := utils.GetCommandExitCode(err)
		if exitcode == 127 {
			fmt.Printf("ERROR: %s\n", "No git executable found")
			os.Exit(exitcode)
		} else if exitcode == 128 {
			fmt.Printf("ERROR: %s\n", "Not a git repository")
			os.Exit(exitcode)
		} else {
			fmt.Println(out)
			os.Exit(exitcode)
		}
	}

	// get git describe but only on semver tags
	cmd = exec.Command("sh", "-c", "git describe --tags --match v*.*.*")
	out, err = cmd.Output()
	if err != nil {
		exitcode := utils.GetCommandExitCode(err)
		if exitcode == 128 {
			// No version tag found, generate initial version
			return generateInitialVersion()
		} else {
			fmt.Println(out)
			os.Exit(exitcode)
		}
	}

	version := strings.TrimSpace(string(out))
	return version
}

func generateInitialVersion() string {
	// Get number of commits
	out, err := exec.Command("sh", "-c", "git rev-list HEAD --count").Output()
	if err != nil {
		fmt.Println(out)
		os.Exit(1)
	}
	num_commits := strings.TrimSpace(string(out))

	// Get short git sha
	out, err = exec.Command("sh", "-c", "git rev-parse --short HEAD").Output()
	if err != nil {
		fmt.Println(out)
		os.Exit(1)
	}
	git_short := strings.TrimSpace(string(out))

	// Create a git describe like snapshot version
	var version = fmt.Sprintf("v0.0.0-%s-g%s", num_commits, git_short)
	return version
}
