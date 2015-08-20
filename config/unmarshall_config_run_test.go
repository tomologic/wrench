package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type UnmarshalConfigRunTestSuite struct {
	suite.Suite
}

func TestUnmarshalConfigRunTestSuite(t *testing.T) {
	suite.Run(t, new(UnmarshalConfigRunTestSuite))
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunSimple() {
	item := yaml.MapItem{
		Key:   "foobar",
		Value: "bash",
	}
	name, run, err := unmarshallConfigRun(item)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foobar", name)
	assert.Equal(suite.T(), "bash", run.Cmd)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunExpandedWithoutEnv() {
	item := yaml.MapItem{
		Key: "foobar",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key:   "Cmd",
				Value: "echo hello",
			},
		},
	}
	name, run, err := unmarshallConfigRun(item)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foobar", name)
	assert.Equal(suite.T(), "echo hello", run.Cmd)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunExpandedWithEnv() {
	env_interface := make([]interface{}, 2)
	env_interface[0] = "FOO=BAR"
	env_interface[1] = "HELLO=WORLD"

	item := yaml.MapItem{
		Key: "foobar",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key:   "Cmd",
				Value: "echo hello",
			},
			yaml.MapItem{
				Key:   "Env",
				Value: env_interface,
			},
		},
	}
	name, run, err := unmarshallConfigRun(item)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foobar", name)
	assert.Equal(suite.T(), "echo hello", run.Cmd)
	assert.Equal(suite.T(), "FOO=BAR", run.Env[0])
	assert.Equal(suite.T(), "HELLO=WORLD", run.Env[1])
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunUnexpectedStructure() {
	item := yaml.MapItem{
		Key: yaml.MapSlice{},
	}
	name, _, err := unmarshallConfigRun(item)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to unmarshall run item", err.Error())
	}
	assert.Equal(suite.T(), "", name)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunEnvUnexpectedStructure() {
	item := yaml.MapItem{
		Key: "foobar",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key:   "Cmd",
				Value: "echo hello",
			},
			yaml.MapItem{
				Key:   "Env",
				Value: "Not a list of interfaces",
			},
		},
	}
	name, _, err := unmarshallConfigRun(item)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to parse Env as list for run item foobar", err.Error())
	}
	assert.Equal(suite.T(), "foobar", name)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunUnknownKey() {
	item := yaml.MapItem{
		Key: "foobar",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key:   "Cmd",
				Value: "echo hello",
			},
			yaml.MapItem{
				Key:   "METADATA",
				Value: "Ignore unknown keys",
			},
		},
	}
	name, run, err := unmarshallConfigRun(item)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foobar", name)
	assert.Equal(suite.T(), "echo hello", run.Cmd)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunSimpleCmdEmpty() {
	item := yaml.MapItem{
		Key:   "simple",
		Value: "      ",
	}
	name, _, err := unmarshallConfigRun(item)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Cmd empty for simple", err.Error())
	}
	assert.Equal(suite.T(), "simple", name)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunExpandedCmdEmpty() {
	item := yaml.MapItem{
		Key: "expanded",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key:   "Cmd",
				Value: "   ",
			},
		},
	}
	name, _, err := unmarshallConfigRun(item)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Cmd empty for expanded", err.Error())
	}
	assert.Equal(suite.T(), "expanded", name)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunExpandedUnexpectedStructure() {
	interface_list := make([]interface{}, 1)
	interface_list[0] = "arbitrary list"

	item := yaml.MapItem{
		Key:   "expanded",
		Value: interface_list,
	}
	name, _, err := unmarshallConfigRun(item)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to parse run item as map for expanded", err.Error())
	}
	assert.Equal(suite.T(), "expanded", name)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunCmdValueNotString() {
	interface_list := make([]interface{}, 1)
	interface_list[0] = "arbitrary list"

	item := yaml.MapItem{
		Key: "expanded",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key:   "Cmd",
				Value: interface_list,
			},
		},
	}
	name, _, err := unmarshallConfigRun(item)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to parse Cmd item as string for run item expanded", err.Error())
	}
	assert.Equal(suite.T(), "expanded", name)
}

func (suite *UnmarshalConfigRunTestSuite) TestUnmarshallConfigRunEnvValueNotString() {
	interface_list := make([]interface{}, 2)
	interface_list[0] = "normal string"
	interface_list[1] = 1

	item := yaml.MapItem{
		Key: "expanded",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key:   "Cmd",
				Value: "run command",
			},
			yaml.MapItem{
				Key:   "Env",
				Value: interface_list,
			},
		},
	}
	name, _, err := unmarshallConfigRun(item)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to parse Env item as string for run item expanded", err.Error())
	}
	assert.Equal(suite.T(), "expanded", name)
}
