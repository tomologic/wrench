package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UnmarshalConfigTestSuite struct {
	suite.Suite
}

func TestUnmarshalConfigTestSuite(t *testing.T) {
	suite.Run(t, new(UnmarshalConfigTestSuite))
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigSimple() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  foobar: bash\n"

	config, err := unmarshallConfig(content)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foobar", config.Project.Name)
	assert.Equal(suite.T(), "bash", config.Run["foobar"].Cmd)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigExpandedWithoutEnv() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  foobar:\n" +
		"    Cmd: bash\n"

	config, err := unmarshallConfig(content)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foobar", config.Project.Name)
	assert.Equal(suite.T(), "bash", config.Run["foobar"].Cmd)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigExpandedWithEnv() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  foobar:\n" +
		"    Cmd: bash\n" +
		"    Env:\n" +
		"      - foobar=1234\n"

	config, err := unmarshallConfig(content)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foobar", config.Project.Name)
	assert.Equal(suite.T(), "bash", config.Run["foobar"].Cmd)
	assert.Equal(suite.T(), "foobar=1234", config.Run["foobar"].Env[0])
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigUnexpectedStructure() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"    - bash\n" +
		"    - foobar=1234\n"

	config, err := unmarshallConfig(content)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to unmarshall Run as map", err.Error())
	}
	assert.Equal(suite.T(), "", config.Project.Name)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigUnexpectedStructureRun() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  true:\n" +
		"    Cmd: bash\n" +
		"    Env:\n" +
		"      - foobar=1234\n"

	config, err := unmarshallConfig(content)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to unmarshall run item", err.Error())
	}
	assert.Equal(suite.T(), "foobar", config.Project.Name)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigEnvUnexpectedStructure() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  true:\n" +
		"    Cmd: bash\n" +
		"    Env: not a list of env vars\n"

	config, err := unmarshallConfig(content)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to unmarshall run item", err.Error())
	}
	assert.Equal(suite.T(), "foobar", config.Project.Name)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigUnknownKey() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  foobar:\n" +
		"    Cmd: bash\n" +
		"    Env:\n" +
		"      - foobar=1234\n" +
		"    Metadata:\n" +
		"      - foobar=1234\n"

	config, err := unmarshallConfig(content)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foobar", config.Project.Name)
	assert.Equal(suite.T(), "bash", config.Run["foobar"].Cmd)
	assert.Equal(suite.T(), "foobar=1234", config.Run["foobar"].Env[0])
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigSimpleCmdEmpty() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  simple: ''\n"

	config, err := unmarshallConfig(content)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Cmd empty for simple", err.Error())
	}
	assert.Equal(suite.T(), "foobar", config.Project.Name)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigExpandedCmdEmpty() {
	content := "Project:\n" +
		"  Name: foobar\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  expanded:\n" +
		"    Cmd: ''\n" +
		"    Env:\n" +
		"      - foobar=1234\n"

	config, err := unmarshallConfig(content)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Cmd empty for expanded", err.Error())
	}
	assert.Equal(suite.T(), "foobar", config.Project.Name)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigExpandedUnexpectedStructure() {
	content := "Project:\n" +
		"  Name: expanded\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  expanded:\n" +
		"    - foobar=1234\n"

	config, err := unmarshallConfig(content)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to parse run item as map for expanded", err.Error())
	}
	assert.Equal(suite.T(), "expanded", config.Project.Name)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigCmdValueNotString() {
	content := "Project:\n" +
		"  Name: expanded\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  expanded:\n" +
		"    Cmd: true\n"

	config, err := unmarshallConfig(content)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to parse Cmd item as string for run item expanded", err.Error())
	}
	assert.Equal(suite.T(), "expanded", config.Project.Name)
}

func (suite *UnmarshalConfigTestSuite) TestUnmarshallConfigEnvValueNotString() {
	content := "Project:\n" +
		"  Name: expanded\n" +
		"  Organization: example\n" +
		"  Version: v1.0.0\n" +
		"Run:\n" +
		"  expanded:\n" +
		"    Cmd: bash\n" +
		"    Env:\n" +
		"      - true\n"

	config, err := unmarshallConfig(content)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Unable to parse Env item as string for run item expanded", err.Error())
	}
	assert.Equal(suite.T(), "expanded", config.Project.Name)
}
