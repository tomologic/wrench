package bump

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/tomologic/wrench/config"
	"github.com/tomologic/wrench/semver"
	"github.com/tomologic/wrench/utils"

	"github.com/spf13/cobra"
)

func AddToWrench(rootCmd *cobra.Command) {
	var cmdBump = &cobra.Command{
		Use:   "bump [major,minor,patch]",
		Short: "Bump project version",
		Long:  `will bump project version, tag git tree and tag snapshot docker image`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) == 0 {
				err = bump("minor")
			} else if len(args) == 1 {
				err = bump(args[0])
			} else {
				cmd.Usage()
				os.Exit(1)
			}

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(cmdBump)
}

func bump(level string) error {
	version, err := semver.Parse(config.GetProjectVersion())
	if err != nil {
		return err
	}

	// Make sure version is snapshot version
	if version.IsReleaseVersion() {
		fmt.Printf("Revision already release '%s'. Doing nothing.\n", version.String())
		os.Exit(0)
	}

	// Create new release version
	if err = version.Bump(level); err != nil {
		return err
	}

	// Make sure docker image of current snapshot version exists
	image_name, err := getImageName()
	if err != nil {
		return err
	}

	// create git tag
	if exitcode, out := utils.RunCmd(fmt.Sprintf("git tag -a %s -m 'Release %s'", version.String(), version.String())); exitcode != 0 {
		return errors.New(fmt.Sprintf("git tag exited with %d: %s\n", exitcode, out))
	}

	// create image
	new_image_name := fmt.Sprintf("%s/%s:%s", config.GetProjectOrganization(), config.GetProjectName(), version.String())
	exitcode, out := utils.RunCmd(fmt.Sprintf("docker tag %s %s", image_name, new_image_name))

	if exitcode != 0 {
		return errors.New(fmt.Sprintf("docker tag exited with %d: %s\n", exitcode, out))
	}

	ver := strings.TrimLeft(version.String(), "v")
	if err := utils.DockerImageAddEnv(new_image_name, "VERSION", ver); err != nil {
		// remove image which is unfinished
		utils.DockerRemoveImage(new_image_name)

		return errors.New("Failed updating VERSION env")
	}

	if exitcode, out := utils.RunCmd(fmt.Sprintf("git push origin %s", version.String())); exitcode != 0 {
		fmt.Printf("ERROR: git push tag exited with %d: %s\n", exitcode, out)
		// remove image
		if exitcode, out := utils.RunCmd(fmt.Sprintf("docker rmi %s", new_image_name)); exitcode != 0 {
			fmt.Printf("ERROR: docker rmi exited with %d: %s\n", exitcode, out)
		}

		// remove git tag
		if exitcode, out := utils.RunCmd(fmt.Sprintf("git tag -d %s", version.String())); exitcode != 0 {
			fmt.Printf("ERROR: git tag delete exited with %d: %s\n", exitcode, out)
		}

		return errors.New("Failed to push new git tag to origin")
	}

	fmt.Printf("Released %s\n", version.String())

	return nil
}

func getImageName() (string, error) {
	git_short, err := getGitShortSha()
	if err != nil {
		return "", err
	}

	tags, err := getGitSemverTags()
	if err != nil {
		return "", err
	}

	// Sort tags with semver lib
	var versions semver.SemverList
	for _, tag := range tags {
		if ver, err := semver.Parse(tag); err != nil {
			return "", err
		} else {
			versions = append(versions, ver)
		}
	}
	sort.Sort(versions)

	// Iterate over all tags
	for _, version := range versions {
		// Convert version back to string
		tag := version.String()

		// Calculate commit count since tag
		num_commits, err := getGitCommitCountSince(tag)
		if err != nil {
			return "", err
		}

		// if count equals 0 then it's not an ancestor
		if num_commits == 0 {
			continue
		}

		// generate snapshot version
		version := fmt.Sprintf("%s-%d-g%s", tag, num_commits, git_short)

		// generate image name
		image_name := fmt.Sprintf("%s/%s:%s",
			config.GetProjectOrganization(),
			config.GetProjectName(),
			version)

		// check if image for this snapshot version exists
		if utils.DockerImageExists(image_name) {
			return image_name, nil
		}
	}

	// check if image for this snapshot version exists based on a root commit
	roots, err := getRootCommits()
	if err != nil {
		return "", err
	}
	for _, sha := range roots {
		// Calculate commit count since root
		num_commits, err := getGitCommitCountSince(sha)
		if err != nil {
			return "", err
		}

		// generate snapshot version
		version := fmt.Sprintf("v0.0.0-%d-g%s", num_commits, git_short)

		// generate image name
		image_name := fmt.Sprintf("%s/%s:%s",
			config.GetProjectOrganization(),
			config.GetProjectName(),
			version)

		// check if image for this snapshot version exists
		if utils.DockerImageExists(image_name) {
			return image_name, nil
		}
	}

	// return error incase image was not found for current revision
	return "", errors.New(fmt.Sprintf("Docker image for revision %s could not be found", git_short))
}

func getGitSemverTags() ([]string, error) {
	exitcode, out := utils.RunCmd("git tag -l 'v[0-9]*\\.[0-9]*\\.[0-9]*'")
	if exitcode != 0 {
		return nil, errors.New(fmt.Sprintf("%d: %s", exitcode, out))
	}

	// Split lines into slice
	tags := strings.Split(out, "\n")

	// Remove empty rows
	tags = utils.RemoveEmptyStrings(tags)

	return tags, nil
}

func getRootCommits() ([]string, error) {
	exitcode, out := utils.RunCmd("git rev-list --max-parents=0 HEAD")
	if exitcode != 0 {
		return nil, errors.New(fmt.Sprintf("%d: %s", exitcode, out))
	}

	// Split lines into slice
	roots := strings.Split(out, "\n")

	// Remove empty rows
	roots = utils.RemoveEmptyStrings(roots)

	return roots, nil
}

func getGitCommitCountSince(sha string) (int, error) {
	exitcode, out := utils.RunCmd(fmt.Sprintf("git rev-list %s..HEAD --count", sha))
	if exitcode != 0 {
		return 0, errors.New(out)
	}

	num, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, err
	}

	return num, nil
}

func getGitShortSha() (string, error) {
	exitcode, out := utils.RunCmd("git rev-parse --short HEAD")
	if exitcode == 128 {
		return "", errors.New("No semver formatted git tag found")
	} else if exitcode != 0 {
		return "", errors.New(out)
	} else if out == "" {
		return "", errors.New("Empty output from git rev-parse")
	}
	return strings.TrimSpace(string(out)), nil
}
