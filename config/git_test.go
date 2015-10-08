package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GitTestSuite struct {
	suite.Suite
}

func TestGitTestSuite(t *testing.T) {
	suite.Run(t, new(GitTestSuite))
}

func (suite *GitTestSuite) TearDownTest() {
	runCmd = mocked_functions["runCmd"].(func(string) (int, string))
}

func (suite *GitTestSuite) TestGitCommitCount() {
	for _, i := range []int{0, 1, 5, 10, 50, 99, 100, 1000, 10000} {
		runCmd = func(string) (int, string) {
			// return +1 since git cli returns total number of commits
			return 0, fmt.Sprintf("%d", (i + 1))
		}

		num_commits, err := getGitCommitCount()

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), i, num_commits)
	}
}

func (suite *GitTestSuite) TestGitCommitUnexpectedString() {
	runCmd = func(string) (int, string) {
		return 0, "foobar"
	}

	num_commits, err := getGitCommitCount()

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 0, num_commits)
}

func (suite *GitTestSuite) TestGitCommitExitCode() {
	runCmd = func(string) (int, string) {
		return 128, "FATAL: unexpected error"
	}

	num_commits, err := getGitCommitCount()

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 0, num_commits)
}

func (suite *GitTestSuite) TestGitShortSha() {
	runCmd = func(string) (int, string) {
		return 0, "aoeu123"
	}

	gitsha, err := getGitShortSha()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "aoeu123", gitsha)
}

func (suite *GitTestSuite) TestGitShortShaEmptyString() {
	runCmd = func(string) (int, string) {
		return 0, ""
	}

	gitsha, err := getGitShortSha()

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", gitsha)
}

func (suite *GitTestSuite) TestGitShortShaExitCode128() {
	runCmd = func(string) (int, string) {
		return 128, "no git semver tag found"
	}

	gitsha, err := getGitShortSha()

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", gitsha)
}

func (suite *GitTestSuite) TestGitShortShaExitCodeUnspecific() {
	runCmd = func(string) (int, string) {
		return 2, "generic exit code"
	}

	gitsha, err := getGitShortSha()

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", gitsha)
}

func (suite *GitTestSuite) TestGitSemverTag() {
	runCmd = func(string) (int, string) {
		return 0, "v123.456.789"
	}

	version, err := getGitSemverTag()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "v123.456.789", version)
}

func (suite *GitTestSuite) TestGitSemverTagEmptyString() {
	runCmd = func(string) (int, string) {
		return 0, ""
	}

	version, err := getGitSemverTag()

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", version)
}

func (suite *GitTestSuite) TestGitSemverTagExitCode128() {
	runCmd = func(string) (int, string) {
		return 128, "no git semver tag found"
	}

	version, err := getGitSemverTag()

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", version)
}

func (suite *GitTestSuite) TestGitSemverTagExitCodeUnspecific() {
	runCmd = func(string) (int, string) {
		return 2, "generic exit code"
	}

	version, err := getGitSemverTag()

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", version)
}
