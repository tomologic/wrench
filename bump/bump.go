package bump

import (
	"fmt"
	"os"
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
	image_name := config.GetProjectImage()

	version, err := semver.Parse(config.GetProjectVersion())
	if err != nil {
		return err
	}

	// Make sure version is snapshot version
	if version.IsReleaseVersion() {
		fmt.Printf("Revision already release '%s'. Doing nothing.\n", version.String())
		os.Exit(0)
	}

	// Make sure docker image of current snapshot version exists
	if !utils.DockerImageExists(image_name) {
		fmt.Printf("ERROR: Docker image %s does not exists\n", image_name)
		os.Exit(1)
	}

	if err = version.Bump(level); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	// create git tag
	if exitcode, out := utils.RunCmd(fmt.Sprintf("git tag -a %s -m 'Release %s'", version.String(), version.String())); exitcode != 0 {
		fmt.Printf("ERROR: git tag exited with %d: %s\n", exitcode, out)
		os.Exit(1)
	}

	// create image
	new_image_name := fmt.Sprintf("%s/%s:%s", config.GetProjectOrganization(), config.GetProjectName(), version.String())
	exitcode, out := utils.RunCmd(
		fmt.Sprintf(
			"docker tag %s %s",
			config.GetProjectImage(),
			new_image_name))

	if exitcode != 0 {
		fmt.Printf("ERROR: docker tag exited with %d: %s\n", exitcode, out)
		os.Exit(1)
	}

	ver := strings.TrimLeft(version.String(), "v")
	if err := utils.DockerImageAddEnv(new_image_name, "VERSION", ver); err != nil {
		// remove image which is unfinished
		utils.DockerRemoveImage(new_image_name)

		fmt.Printf("ERROR: Failed updating VERSION env\n")
		os.Exit(1)
	}

	// push tag; if err
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

		os.Exit(1)
	}

	fmt.Printf("Released %s\n", version.String())

	return nil
}
