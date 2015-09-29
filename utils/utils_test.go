package utils

import (
	"archive/tar"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestRunCmd() {
	exitcode, out := RunCmd("echo hello")

	if assert.Equal(suite.T(), 0, exitcode) {
		assert.Equal(suite.T(), "hello\n", out)
	}
}

func (suite *UtilsTestSuite) TestRunExitCode1() {
	exitcode, out := RunCmd("exit 1")

	if assert.Equal(suite.T(), 1, exitcode) {
		assert.Equal(suite.T(), "", out)
	}
}

func (suite *UtilsTestSuite) TestRunExitCodeCustom() {
	exitcode, out := RunCmd("bash -c 'echo -n test foo bar' && exit 127")

	if assert.Equal(suite.T(), 127, exitcode) {
		assert.Equal(suite.T(), "test foo bar", out)
	}
}

func (suite *UtilsTestSuite) TestCreateTar() {
	var files = make([]Tarfile, 2)
	files[0] = Tarfile{
		"readme.txt",
		"This archive contains some text files.",
	}
	files[1] = Tarfile{
		"gopher.txt",
		"Gopher names:\nGeorge\nGeoffrey\nGonzo",
	}

	tarfile, err := CreateTar(files)
	if assert.Nil(suite.T(), err) {

		// Open the tar archive for reading.
		r := bytes.NewReader(tarfile.Bytes())
		tr := tar.NewReader(r)

		// Iterate through all files
		for i := 0; i < len(files); i++ {
			// Get next file in tar
			hdr, err := tr.Next()
			if assert.Nil(suite.T(), err) {
				// Check filename
				assert.Equal(suite.T(), files[i].Name, hdr.Name)

				content := new(bytes.Buffer)
				content.ReadFrom(tr)

				// Check content
				assert.Equal(suite.T(), files[i].Content, content.String())
			}
		}
	}
}
