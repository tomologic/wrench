package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
	Cmd     string   `yaml:"Cmd"`
	Env     []string `yaml:"Env,omitempty"`
	Volumes []string `yaml:"Volume,omitempty"`
}
type Config struct {
	Project Project        `yaml:"Project"`
	Run     map[string]Run `yaml:"Run,omitempty"`
}
type TemplateContext struct {
	Environ *map[string]string
}

var config = &Config{}

func AddToWrench(cmdRoot *cobra.Command) {
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
	c, err := loadConfigFile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config = &c
}

func commandConfig(format string) string {
	// Detect values for fields not set
	GetProjectOrganization()
	GetProjectName()
	GetProjectVersion()
	GetProjectImage()

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

var getEnviron = func() []string {
	return os.Environ()
}

var getTmplContextEnviron = func() *map[string]string {
	Environ := make(map[string]string)

	// Get all environment variables
	for _, item := range getEnviron() {
		splits := strings.Split(item, "=")
		Environ[splits[0]] = strings.Join(splits[1:], "=")
	}

	return &Environ
}

var getConfigContent = func() (string, error) {
	if !utils.FileExists("./wrench.yml") {
		return "", nil
	}

	content, err := ioutil.ReadFile("./wrench.yml")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

var getRenderedConfigContent = func(content string) (string, error) {
	// Create template from wrench file
	var config_rendered bytes.Buffer
	tmpl, err := template.New("config").Parse(content)
	if err != nil {
		return "", err
	}

	// Get context for template
	tmpl_context := TemplateContext{
		Environ: getTmplContextEnviron(),
	}

	// Render template with tmpl_context
	err = tmpl.Execute(&config_rendered, tmpl_context)
	if err != nil {
		return "", err
	}

	return config_rendered.String(), nil
}

var unmarshallConfigRun = func(item yaml.MapItem) (string, Run, error) {
	run := Run{}

	name, ok := item.Key.(string)
	if ok != true {
		return name, run, errors.New("Unable to unmarshall run item")
	}

	// No more values provided if value is a string
	cmd_string, ok := item.Value.(string)
	if ok {
		// Simple string structure
		// foobar: run command
		run.Cmd = strings.TrimSpace(cmd_string)
	} else {
		// Expanded structure
		// foobar:
		//   Cmd: run command
		//   Env:
		//     - ENVVAR=value
		//   Volumes:
		//     - /localdir:/containerdir
		run_expanded, ok := item.Value.(yaml.MapSlice)
		if !ok {
			return name, run, errors.New(fmt.Sprintf("Unable to parse run item as map for %s", name))
		}

		for k := range run_expanded {
			key := run_expanded[k].Key.(string)
			if key == "Cmd" {
				// Expecting a command string
				cmd_string, ok = run_expanded[k].Value.(string)
				if !ok {
					return name, run, errors.New(fmt.Sprintf("Unable to parse Cmd item as string for run item %s", name))
				}
				run.Cmd = strings.TrimSpace(cmd_string)
			} else if key == "Env" || key == "Volumes" {
				// Expecting a list of items
				item_list, ok := run_expanded[k].Value.([]interface{})
				if !ok {
					return name, run, errors.New(fmt.Sprintf("Unable to parse %s as list for run item %s", key, name))
				}
				for _, s := range item_list {
					t, ok := s.(string)
					if ok {
						if key == "Env" {
							run.Env = append(run.Env, t)
						} else if key == "Volumes" {
							run.Volumes = append(run.Volumes, t)
						}
					} else {
						return name, run, errors.New(fmt.Sprintf("Unable to parse %s item as string for run item %s", key, name))
					}
				}
			}
		}
	}

	if run.Cmd == "" {
		return name, run, errors.New(fmt.Sprintf("Cmd empty for %s", name))
	}

	return name, run, nil
}

var unmarshallConfig = func(content string) (Config, error) {
	type UnmarshalConfig struct {
		Project Project       `yaml:"Project"`
		Run     yaml.MapSlice `yaml:"Run,omitempty"`
	}

	uconfig := UnmarshalConfig{}
	config := Config{}

	// Load the expected yaml file structure
	err := yaml.Unmarshal([]byte(content), &uconfig)
	if err != nil {
		return config, errors.New("Unable to unmarshall Run as map")
	}

	// Get Project from unmarshalled config
	config.Project = uconfig.Project

	// Create Run map in config
	config.Run = make(map[string]Run)

	// Unmarshal every run item in Run map
	for _, item := range uconfig.Run {
		name, run, err := unmarshallConfigRun(item)
		if err != nil {
			return config, err
		}
		config.Run[name] = run
	}

	return config, nil
}

func loadConfigFile() (Config, error) {
	// Get wrench file content
	config_content, err := getConfigContent()

	if err != nil {
		return Config{}, err
	}

	// Return if config file content is empty
	if config_content == "" {
		return Config{}, nil
	}

	config_rendered, err := getRenderedConfigContent(config_content)
	if err != nil {
		return Config{}, err
	}

	// Return if rendered_config content is empty
	if config_rendered == "" {
		return Config{}, nil
	}

	config, err := unmarshallConfig(config_rendered)
	if err != nil {
		return Config{}, err
	}

	return config, nil
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

var getHostname = func() (string, error) {
	exitcode, out := runCmd("hostname -f")
	if exitcode != 0 {
		return "", errors.New(fmt.Sprintf("hostname exited with %d", exitcode))
	}
	return out, nil
}

func detectProjectOrganization() string {
	hostname, err := getHostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	parts := strings.Split(string(hostname), ".")

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
	return utils.RunCmd(command)
}

var getGitRepoPresent = func() (bool, error) {
	exitcode, out := runCmd("git rev-parse --short HEAD")
	if exitcode == 127 {
		return false, errors.New("No git executable found")
	} else if exitcode == 128 {
		return false, errors.New("Not a git repository")
	} else if exitcode != 0 {
		return false, errors.New(out)
	}
	return true, nil
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
	if present, err := getGitRepoPresent(); !present {
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

	num, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, err
	}

	// Get number of commits since initial commit instead of total
	num -= 1

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
