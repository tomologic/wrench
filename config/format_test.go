package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FormatTestSuite struct {
	suite.Suite
}

func TestFormatTestSuite(t *testing.T) {
	suite.Run(t, new(FormatTestSuite))
}

func (suite *FormatTestSuite) SetupSuite() {
	// Setup prefined data in config for format test
	config = &Config{
		Project: Project{
			Organization: "example",
			Name:         "foobar",
			Version:      "v1.0.0",
		},
		Run: map[string]Run{
			"syntax-test": {Cmd: "flake8 -v ."},
		},
	}
}

func (suite *FormatTestSuite) TearDownSuite() {
	// Cleanup state after suite
	config = &Config{}
}

func (suite *FormatTestSuite) TestConfigCommand() {
	expected := "Project:\n" +
		"  Organization: example\n" +
		"  Name: foobar\n" +
		"  Version: v1.0.0\n" +
		"  Image: example/foobar:v1.0.0\n" +
		"Run:\n" +
		"  syntax-test:\n" +
		"    Cmd: flake8 -v .\n"
	assert.Equal(suite.T(), expected, commandConfig(""))
}

func (suite *FormatTestSuite) TestConfigCommandFormatProject() {
	assert.Equal(suite.T(), "{example foobar v1.0.0 example/foobar:v1.0.0}", commandConfig("{{.Project}}"))
}

func (suite *FormatTestSuite) TestConfigCommandFormatProjectName() {
	assert.Equal(suite.T(), "foobar", commandConfig("{{.Project.Name}}"))
}

func (suite *FormatTestSuite) TestConfigCommandFormatProjectVersion() {
	assert.Equal(suite.T(), "v1.0.0", commandConfig("{{.Project.Version}}"))
}

func (suite *FormatTestSuite) TestConfigCommandFormatProjectImage() {
	assert.Equal(suite.T(), "example/foobar:v1.0.0", commandConfig("{{.Project.Image}}"))
}
