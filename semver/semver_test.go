package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SemverTestSuite struct {
	suite.Suite
}

func TestSemverTestSuite(t *testing.T) {
	suite.Run(t, new(SemverTestSuite))
}

func (suite *SemverTestSuite) TestSemverPartEmpty() {
	var examples = []struct {
		Input string
		Error string
	}{
		{"..", "Part 1 empty"},
		{".0.1", "Part 1 empty"},
		{"1..1", "Part 2 empty"},
		{"1.2.", "Part 3 empty"},
		{"v..", "Part 1 empty"},
		{"v.0.1", "Part 1 empty"},
		{"v1..1", "Part 2 empty"},
		{"v1.2.", "Part 3 empty"},
	}

	for _, ex := range examples {
		_, err := Parse(ex.Input)

		if assert.NotNil(suite.T(), err) {
			assert.Equal(suite.T(), ex.Error, err.Error())
		}
	}
}

func (suite *SemverTestSuite) TestSemverRequire3Parts() {
	examples := []string{
		"1",
		"v1",
		"v1-1-gaoeu123",
		"v1.2-1-gaoeu123",
	}
	for _, ex := range examples {
		_, err := Parse(ex)

		if assert.NotNil(suite.T(), err) {
			assert.Equal(suite.T(), "Semver requires 3 parts Major.Minor.Patch[-snapshot]", err.Error())
		}
	}
}

func (suite *SemverTestSuite) TestSemver() {
	var examples = []struct {
		Input    string
		Expected Semver
	}{
		{"1.2.3", Semver{1, 2, 3, ""}},
		{"123.456.789", Semver{123, 456, 789, ""}},
		{"1.2.3-1-gaoeu123", Semver{1, 2, 3, "1-gaoeu123"}},
		{"123.456.789-1-gaoeu123", Semver{123, 456, 789, "1-gaoeu123"}},
		{"v1.2.3", Semver{1, 2, 3, ""}},
		{"v123.456.789", Semver{123, 456, 789, ""}},
		{"v1.2.3-1-gaoeu123", Semver{1, 2, 3, "1-gaoeu123"}},
		{"v123.456.789-1-gaoeu123", Semver{123, 456, 789, "1-gaoeu123"}},
	}

	for _, ex := range examples {
		semver, err := Parse(ex.Input)

		assert.Nil(suite.T(), err)
		assert.True(suite.T(), semver.Equals(ex.Expected))
	}
}

func (suite *SemverTestSuite) TestSemverBumpUnknownLevel() {
	sv := Semver{1, 2, 3, ""}
	err := sv.Bump("")

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "Unknown level ''", err.Error())

	err = sv.Bump("aoeu")

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "Unknown level 'aoeu'", err.Error())
}

func (suite *SemverTestSuite) TestSemverBump() {
	var examples = []struct {
		level    string
		Input    Semver
		Expected Semver
	}{
		{
			"Major",
			Semver{1, 2, 3, ""},
			Semver{2, 0, 0, ""},
		},
		{
			"major",
			Semver{1, 2, 3, "1-g1234567"},
			Semver{2, 0, 0, ""},
		},
		{
			"major",
			Semver{1, 2, 3, "snapshot"},
			Semver{2, 0, 0, ""},
		},
		{
			"major",
			Semver{123, 456, 789, "1-gaoeu123"},
			Semver{124, 0, 0, ""},
		},
		{
			"Minor",
			Semver{1, 2, 3, ""},
			Semver{1, 3, 0, ""},
		},
		{
			"minor",
			Semver{1, 2, 3, "1-g1234567"},
			Semver{1, 3, 0, ""},
		},
		{
			"minor",
			Semver{1, 2, 3, "snapshot"},
			Semver{1, 3, 0, ""},
		},
		{
			"minor",
			Semver{123, 456, 789, "1-gaoeu123"},
			Semver{123, 457, 0, ""},
		},
		{
			"Patch",
			Semver{1, 2, 3, ""},
			Semver{1, 2, 4, ""},
		},
		{
			"patch",
			Semver{1, 2, 3, "1-g1234567"},
			Semver{1, 2, 4, ""},
		},
		{
			"patch",
			Semver{1, 2, 3, "snapshot"},
			Semver{1, 2, 4, ""},
		},
		{
			"patch",
			Semver{123, 456, 789, "1-gaoeu123"},
			Semver{123, 456, 790, ""},
		},
	}

	for _, ex := range examples {
		err := ex.Input.Bump(ex.level)

		if assert.Nil(suite.T(), err) {
			assert.Equal(suite.T(), ex.Expected.String(), ex.Input.String())
		}
	}
}

func (suite *SemverTestSuite) TestSemverString() {
	var examples = []struct {
		Expected string
		Input    Semver
	}{
		{"v1.2.3", Semver{1, 2, 3, ""}},
		{"v123.456.789", Semver{123, 456, 789, ""}},
		{"v1.2.3-1-gaoeu123", Semver{1, 2, 3, "1-gaoeu123"}},
		{"v123.456.789-1-gaoeu123", Semver{123, 456, 789, "1-gaoeu123"}},
	}

	for _, ex := range examples {
		str := ex.Input.String()

		assert.Equal(suite.T(), ex.Expected, str)
	}
}
