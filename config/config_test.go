package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	config = &Config{
		Project: Project{
			Organization: "example",
			Name:         "foobar",
			Version:      "v1.0.0",
		},
	}

	assert.Equal(t, GetProjectOrganization(), "example")
	assert.Equal(t, GetProjectName(), "foobar")
	assert.Equal(t, GetProjectVersion(), "v1.0.0")
}

func ExampleConfigCommand() {
	config = &Config{
		Project: Project{
			Organization: "example",
			Name:         "foobar",
			Version:      "v1.0.0",
		},
	}

	commandConfig()
	// Output:
	// Project:
	//   Organization: example
	//   Name: foobar
	//   Version: v1.0.0
}

func ExampleConfigCommandWithRun() {
	config = &Config{
		Project: Project{
			Organization: "example",
			Name:         "foobar",
			Version:      "v1.0.0",
		},
		Run: map[string]string{
			"syntax-test": "flake8 -v .",
		},
	}

	commandConfig()
	// Output:
	// Project:
	//   Organization: example
	//   Name: foobar
	//   Version: v1.0.0
	// Run:
	//   syntax-test: flake8 -v .
}

func ExampleConfigCommandFormatProject() {
	config = &Config{
		Project: Project{
			Organization: "example",
			Name:         "foobar",
			Version:      "v1.0.0",
		},
	}

	flag_format = "{{.Project}}"
	commandConfig()
	// Output: {example foobar v1.0.0}
}

func ExampleConfigCommandFormatProjectName() {
	config = &Config{
		Project: Project{
			Organization: "example",
			Name:         "foobar",
			Version:      "v1.0.0",
		},
	}

	flag_format = "{{.Project.Name}}"
	commandConfig()
	// Output: foobar
}

func ExampleConfigCommandFormatProjectVersion() {
	config = &Config{
		Project: Project{
			Organization: "example",
			Name:         "foobar",
			Version:      "v1.0.0",
		},
	}

	flag_format = "{{.Project.Version}}"
	commandConfig()
	// Output: v1.0.0
}
