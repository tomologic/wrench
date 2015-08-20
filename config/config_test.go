package config

var mocked_functions = map[string]interface{}{
	"getAbsFilePath":         getAbsFilePath,
	"generateInitialVersion": generateInitialVersion,
	"getHostname":            getHostname,
	"getGitCommitCount":      getGitCommitCount,
	"getGitSemverTag":        getGitSemverTag,
	"getGitShortSha":         getGitShortSha,
	"getGitRepoPresent":      getGitRepoPresent,
	"runCmd":                 runCmd,
	"unmarshallConfigRun":    unmarshallConfigRun,
}
