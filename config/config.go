package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func AddToWrench(cmdRoot *cobra.Command) {
	readWrenchFile()

	var flag_format string

	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Configuration for wrench",
		Long:  `configuration picked up by wrench and used in commands`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(commandConfig(flag_format))
		},
	}

	cmdConfig.Flags().StringVar(&flag_format, "format", "", "Return specific value from config")

	cmdRoot.AddCommand(cmdConfig)
}

func commandConfig(format string) string {
	generateAllConfig()

	if format == "" {
		d, err := yaml.Marshal(&config)
		if err != nil {
			panic(err)
		}
		return string(d)
	} else {
		tmpl, err := template.New("format").Parse(format)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		var out bytes.Buffer
		err = tmpl.Execute(&out, &config)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		return out.String()
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

var getFqdn = func() string {
	fqdn, err := exec.Command("sh", "-c", "hostname -f").Output()
	if err != nil {
		panic(err)
	}
	return string(fqdn)
}

func detectProjectOrganization() string {
	fqdn := getFqdn()
	parts := strings.Split(string(fqdn), ".")

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

var getAbsFilePath = func() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func detectProjectName() string {
	dir := getAbsFilePath()
	if project := string(filepath.Base(dir)); project == "/" {
		return "noname"
	} else {
		return project
	}
}

var runCmd = func(command string) (int, string) {
	exitcode := 0
	cmd := exec.Command(fmt.Sprintf("sh -c %s", command))
	out, err := cmd.Output()
	if err != nil {
		exitcode = utils.GetCommandExitCode(err)
	}
	return exitcode, string(out)
}

var gitRepoPresent = func() error {
	exitcode, out := runCmd("git rev-parse --short HEAD")
	if exitcode == 127 {
		return errors.New("No git executable found")
	} else if exitcode == 128 {
		return errors.New("Not a git repository")
	} else if exitcode != 0 {
		return errors.New(out)
	}
	return nil
}

var getGitSemverTag = func() (string, error) {
	// get git describe but only on semver tags
	exitcode, out := runCmd("git describe --tags --match v*.*.*")
	if exitcode == 128 {
		// No version tag found, generate initial version
		return "", errors.New("No semver formatted git tag found")
	} else if exitcode != 0 {
		return "", errors.New(out)
	} else if out == "" {
		return "", errors.New("Empty output from git describe")
	}

	version := strings.TrimSpace(string(out))
	return version, nil
}

func detectProjectVersion() string {
	// make sure git is installed and we are inside a git repo
	if err := gitRepoPresent(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get latest git semver version
	if version, err := getGitSemverTag(); err != nil {
		return generateInitialVersion()
	} else {
		return version
	}
}

var getGitCommitCount = func() (int, error) {
	exitcode, out := runCmd("git rev-list HEAD --count")
	if exitcode != 0 {
		return 0, errors.New(out)
	}

	num, err := strconv.Atoi(out)
	if err != nil {
		return 0, err
	}

	return num, nil
}

var getGitShortSha = func() (string, error) {
	exitcode, out := runCmd("git rev-parse --short HEAD")
	if exitcode == 128 {
		return "", errors.New("No semver formatted git tag found")
	} else if exitcode != 0 {
		return "", errors.New(out)
	} else if out == "" {
		return "", errors.New("Empty output from git rev-parse")
	}
	return strings.TrimSpace(string(out)), nil
}

var generateInitialVersion = func() string {
	// Get number of commits
	num_commits, err := getGitCommitCount()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get short git sha
	git_short, err := getGitShortSha()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a git describe like snapshot version
	var version = fmt.Sprintf("v0.0.0-%d-g%s", num_commits, git_short)
	return version
}
