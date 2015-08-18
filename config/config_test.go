package config

var mocked_functions = map[string]interface{}{
	"getAbsFilePath":         getAbsFilePath,
	"generateInitialVersion": generateInitialVersion,
	"getFqdn":                getFqdn,
	"getGitCommitCount":      getGitCommitCount,
	"getGitSemverTag":        getGitSemverTag,
	"getGitShortSha":         getGitShortSha,
	"gitRepoPresent":         gitRepoPresent,
	"runCmd":                 runCmd,
}
