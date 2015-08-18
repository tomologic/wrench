package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DetectTestSuite struct {
	suite.Suite
}

func TestDetectTestSuite(t *testing.T) {
	suite.Run(t, new(DetectTestSuite))
}

func (suite *DetectTestSuite) SetupSuite() {
	gitRepoPresent = func() error { return nil }
}

func (suite *DetectTestSuite) TearDownSuite() {
	gitRepoPresent = mocked_functions["gitRepoPresent"].(func() error)
}

func (suite *DetectTestSuite) TearDownTest() {
	getAbsFilePath = mocked_functions["getAbsFilePath"].(func() string)
	getFqdn = mocked_functions["getFqdn"].(func() string)
	runCmd = mocked_functions["runCmd"].(func(string) (int, string))
	getGitSemverTag = mocked_functions["getGitSemverTag"].(func() (string, error))
	generateInitialVersion = mocked_functions["generateInitialVersion"].(func() string)
}

func (suite *DetectTestSuite) TestDetectProjectName() {
	getAbsFilePath = func() string {
		return "/foobar"
	}
	assert.Equal(suite.T(), "foobar", detectProjectName())
}

func (suite *DetectTestSuite) TestDetectProjectNameRoot() {
	getAbsFilePath = func() string {
		return "/"
	}
	assert.Equal(suite.T(), "noname", detectProjectName())
}

func (suite *DetectTestSuite) TestDetectProjectNameSub() {
	getAbsFilePath = func() string {
		return "/home/root/aoeu1234"
	}
	assert.Equal(suite.T(), "aoeu1234", detectProjectName())
}

func (suite *DetectTestSuite) TestDetectProjectOrganizationSingle() {
	getFqdn = func() string {
		return "user"
	}
	assert.Equal(suite.T(), "user", detectProjectOrganization())
}

func (suite *DetectTestSuite) TestDetectProjectOrganizationLocal() {
	getFqdn = func() string {
		return "user.local"
	}
	assert.Equal(suite.T(), "local", detectProjectOrganization())
}

func (suite *DetectTestSuite) TestDetectProjectOrganizationDomain() {
	getFqdn = func() string {
		return "hostname.domain.topdomain"
	}
	assert.Equal(suite.T(), "domain", detectProjectOrganization())
}

func (suite *DetectTestSuite) TestDetectProjectOrganizationSubDomain() {
	getFqdn = func() string {
		return "hostname.subdomain.domain.topdomain"
	}
	assert.Equal(suite.T(), "domain", detectProjectOrganization())
}

func (suite *DetectTestSuite) TestDetectProjectVersion() {
	runCmd = func(string) (int, string) {
		return 0, "v0.1.0-1-g1234567"
	}
	assert.Equal(suite.T(), "v0.1.0-1-g1234567", detectProjectVersion())
}

func (suite *DetectTestSuite) TestDetectProjectVersionRelease() {
	runCmd = func(string) (int, string) {
		return 0, "v123.456.789"
	}
	assert.Equal(suite.T(), "v123.456.789", detectProjectVersion())
}

func (suite *DetectTestSuite) TestDetectProjectVersionNoTag() {
	runCmd = func(string) (int, string) {
		return 128, "No semver formatted git tag found"
	}
	generateInitialVersion = func() string { return "generated-version" }

	assert.Equal(suite.T(), "generated-version", detectProjectVersion())
}
