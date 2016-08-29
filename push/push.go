package push

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tomologic/wrench/config"
	"github.com/tomologic/wrench/utils"

	"github.com/spf13/cobra"
)

func AddToWrench(rootCmd *cobra.Command) {
	var flag_additional_tags string

	var cmdBump = &cobra.Command{
		Use:   "push [--additional-tags]",
		Short: "Push project release image to docker registry",
		Long:  `will push release image for current version to specified registry`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				if err := push(args[0], flag_additional_tags); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				cmd.Usage()
				os.Exit(1)
			}
		},
	}

	cmdBump.Flags().StringVar(&flag_additional_tags, "additional-tags", "", "Comma separated list of additional tags to push 'latest,prod'")

	rootCmd.AddCommand(cmdBump)
}

func push(registry string, additional_tags string) error {
	tags := strings.Split(additional_tags, ",")
	tags = append(tags, config.GetProjectVersion())
	tags = utils.RemoveEmptyStrings(tags)

	for _, tag := range tags {
		image_name := fmt.Sprintf("%s/%s:%s",
			config.GetProjectOrganization(),
			config.GetProjectName(),
			config.GetProjectVersion())

		registry = strings.Trim(registry, "/")

		tmp_image_name := fmt.Sprintf("%s/%s/%s:%s",
			registry,
			config.GetProjectOrganization(),
			config.GetProjectName(),
			tag)

		// prefix image name with registry
		if err := tag_image(image_name, tmp_image_name); err != nil {
			return err
		}

		// push prefixed image
		push_err := push_image(tmp_image_name)

		// cleanup prefixed images
		cleanup_err := remove_image(tmp_image_name)

		if push_err != nil {
			return push_err
		}

		if cleanup_err != nil {
			return cleanup_err
		}
	}

	return nil
}

func tag_image(image_name string, new_image_name string) error {
	command := fmt.Sprintf(
		"docker tag %s %s",
		image_name,
		new_image_name)

	exitcode, out := utils.RunCmd(command)
	if exitcode != 0 {
		fmt.Fprintln(os.Stderr, out)

		return errors.New(fmt.Sprintf(
			"Could not retag %s to %s",
			image_name,
			new_image_name))
	}

	return nil
}

func push_image(image string) error {
	command := fmt.Sprintf(
		"docker push %s", image)

	exitcode, out := utils.RunCmd(command)
	if exitcode != 0 {
		fmt.Fprintln(os.Stderr, out)

		return errors.New(fmt.Sprintf("Could not push %s", image))
	}

	return nil
}

func remove_image(image string) error {
	command := fmt.Sprintf(
		"docker rmi %s", image)

	exitcode, out := utils.RunCmd(command)
	if exitcode != 0 {
		fmt.Fprintln(os.Stderr, out)

		return errors.New(fmt.Sprintf(
			"Could not remove %s", image))
	}

	return nil
}
