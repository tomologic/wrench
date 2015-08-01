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

var run_list *map[string]interface{}

func AddToWrench(cmdRoot *cobra.Command) {
	var cmdRun = &cobra.Command{
		Use:   "run [command]",
		Short: "Run commands in docker image",
		Long:  `Run defined commands from wrench.yml inside application image`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Usage()
				os.Exit(1)
			}

			run(args[0])
		},
	}

	cmdRoot.AddCommand(cmdRun)
}

func run(name string) {
	image_name := config.GetProjectImage()

	run, ok := config.GetRun(name)
	if ok == false {
		fmt.Println("ERROR: %s not found in wrench.yml", name)
		os.Exit(1)
	}

	if utils.FileExists("./Dockerfile.test") {
		// If test dockerfile exists then use test image
		image_name = fmt.Sprintf("%s-test", image_name)
	}

	if !utils.DockerImageExists(image_name) {
		fmt.Printf("ERROR: Image %s does not exist, run wrench build\n", image_name)
		os.Exit(1)
	}

	fmt.Printf("INFO: running %s in image %s\n", name, image_name)

	// Tempdir for building temporary run image
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	tempdir, err := ioutil.TempDir(dir, ".wrench_run_")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create temp Dockerfile
	dockerfile_content := "" +
		fmt.Sprintf("FROM %s\n", image_name) +
		"ADD wrench_run.sh /tmp/\n" +
		"ENTRYPOINT [\"bash\"]\n" +
		"CMD [\"/tmp/wrench_run.sh\"]\n"

	dockerfile := fmt.Sprintf("%s/Dockerfile", tempdir)
	utils.WriteFileContent(dockerfile, dockerfile_content)

	// Create wrench run bash file
	runfile := fmt.Sprintf("%s/wrench_run.sh", tempdir)
	utils.WriteFileContent(runfile, run.Cmd)

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

	// Defer so all defered cleanup is done
	defer func() {
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(utils.GetCommandExitCode(err))
		}
	}()
	defer utils.DockerRemoveImage(run_image_name)
	defer os.RemoveAll(tempdir)

	if len(run.Env) > 0 {
		envfile := fmt.Sprintf("./%s-env", tempdir_base)
		utils.WriteFileContent(envfile, strings.Join(run.Env, "\n"))
		defer os.Remove(envfile)

		cmd_string = fmt.Sprintf("docker run -t --rm --env-file '%s' '%s'", envfile, run_image_name)
	} else {
		cmd_string = fmt.Sprintf("docker run -t --rm '%s'", run_image_name)
	}

	// Run
	cmd = exec.Command("sh", "-c", cmd_string)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
}
