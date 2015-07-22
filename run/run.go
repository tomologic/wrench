package run

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomologic/wrench/config"
	"github.com/tomologic/wrench/utils"
)

var run_list *map[string]string

func AddToWrench(cmdRoot *cobra.Command) {
	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Run commands in docker image",
		Long:  `Run defined commands from wrench.yml inside application image`,
		Run: func(cmd *cobra.Command, args []string) {
			run_list = config.GetRunList()

			run(args[0])
		},
	}

	cmdRoot.AddCommand(cmdRun)
}

func run(name string) {
	image_name := fmt.Sprintf("%s/%s:%s",
		config.GetProjectOrganization(),
		config.GetProjectName(),
		config.GetProjectVersion())

	if utils.FileExists("./Dockerfile.test") {
		// If test dockerfile exists then use test image
		image_name = fmt.Sprintf("%s-test", image_name)
	}

	if !utils.DockerImageExists(image_name) {
		fmt.Printf("ERROR: Image %s does not exist, run wrench build\n", image_name)
		os.Exit(1)
	}

	runfile_content := []string{(*run_list)[name]}
	if strings.TrimSpace(runfile_content[0]) == "" {
		fmt.Printf("ERROR: Command %s is empty\n", name)
		os.Exit(1)
	}

	fmt.Printf("INFO: running %s in image %s\n", name, image_name)

	// Tempdir for building temporary run image
	dir := os.TempDir()
	tempdir, err := ioutil.TempDir(dir, "wrench_run")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempdir)

	// Create temp Dockerfile
	dockerfile_content := []string{
		fmt.Sprintf("FROM %s", image_name),
		"ADD wrench_run.sh /tmp/",
		"ENTRYPOINT [\"bash\"]",
		"CMD [\"/tmp/wrench_run.sh\"]",
	}

	dockerfile := fmt.Sprintf("%s/Dockerfile", tempdir)
	utils.WriteFileContent(dockerfile, dockerfile_content)

	// Create wrench run bash file
	runfile := fmt.Sprintf("%s/wrench_run.sh", tempdir)
	utils.WriteFileContent(runfile, runfile_content)

	tempdir_base := string(filepath.Base(tempdir))
	run_image_name := fmt.Sprintf("%s-%s", image_name, tempdir_base)

	cmd_string := fmt.Sprintf("docker build -t '%s' .", run_image_name)
	cmd := exec.Command("sh", "-c", cmd_string)
	cmd.Dir = tempdir
	out, err := cmd.Output()

	// Remove tempdir
	if err != nil {
		os.RemoveAll(tempdir)

		fmt.Printf("ERROR: %s\n", string(out))
		os.Exit(1)
	}
	defer utils.DockerRemoveImage(run_image_name)

	// Run
	cmd_string = fmt.Sprintf("docker run -t --rm '%s'", run_image_name)
	cmd = exec.Command("sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		os.RemoveAll(tempdir)
		utils.DockerRemoveImage(run_image_name)

		fmt.Println(err)
		os.Exit(utils.GetCommandExitCode(err))
	}

}
