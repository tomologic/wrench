package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InitialVersionTestSuite struct {
	suite.Suite
}

func TestInitialVersionTestSuite(t *testing.T) {
	suite.Run(t, new(InitialVersionTestSuite))
}

func (suite *InitialVersionTestSuite) TearDownTest() {
	getGitShortSha = mocked_functions["getGitShortSha"].(func() (string, error))
	getGitCommitCount = mocked_functions["getGitCommitCount"].(func() (int, error))
}

func (suite *InitialVersionTestSuite) TestInitialVersion() {
	getGitShortSha = func() (string, error) {
		return "aoeu123", nil
	}
	getGitCommitCount = func() (int, error) {
		return 1, nil
	}

	assert.Equal(suite.T(), "v0.0.0-1-gaoeu123", generateInitialVersion())
}

func (suite *InitialVersionTestSuite) TestInitialVersionFoobar() {
	getGitShortSha = func() (string, error) {
		return "foobar", nil
	}
	for _, i := range []int{0, 1, 5, 10, 50, 99, 100, 1000, 10000} {
		getGitCommitCount = func() (int, error) {
			return i, nil
		}

		assert.Equal(suite.T(), fmt.Sprintf("v0.0.0-%d-gfoobar", i), generateInitialVersion())
	}
}
