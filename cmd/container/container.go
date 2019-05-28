package container

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.bearstech.com/factory/factory-cli/cmd/root"
)

var (
	dry_run     bool
	target      string
	environment string
)

func init() {
	root.FlagE(execCmd.PersistentFlags(), &environment)
	execCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")

	containerCmd.AddCommand(execCmd)
	root.RootCmd.AddCommand(containerCmd)
}

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Do something on a container",
	Long:  `Do something on a container.`,
}

var execCmd = &cobra.Command{
	Use:   "exec service [command]",
	Short: "Exec a command in a container",
	Long:  `Exec a command in a container. Default to bash`,
	Args: func(cmd *cobra.Command, args []string) error {
		if root.Project == "" {
			return errors.New("please specify a project with -p")
		}
		if err := root.AssertEnvironment(environment); err != nil {
			return err
		}
		if len(args) == 0 {
			return errors.New("you must specify a service")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := root.SSHAddress(environment)
		if err != nil {
			return err
		}
		command := []string{
			"ssh",
			"-a", "-x", "-t", "-p", "2222",
			address,
			"exec",
		}
		if len(args) < 1 {
			args = []string{args[0], "bash"}
		}
		command = append(command, args...)

		log.Debug(command)
		if dry_run {
			fmt.Printf("%s\n", strings.Join(command, " "))
		} else {
			c := exec.Command("ssh")
			c.Args = command
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			err := c.Run()
			if err != nil {
				return err
			}
		}
		return nil
	},
}
