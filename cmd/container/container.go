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
	"gitlab.bearstech.com/factory/factory-cli/signpost"
)

var (
	dry_run     bool
	target      string
	environment string
)

func init() {
	execCmd.PersistentFlags().StringVarP(&environment, "environment", "E", "", "Environment")
	execCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")

	containerCmd.AddCommand(execCmd)
	root.RootCmd.AddCommand(containerCmd)
}

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Do something on a container",
	Long:  `Do something on a container.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("you must use a subcommand: exec")
	},
}

var execCmd = &cobra.Command{
	Use:   "exec service [command]",
	Short: "Exec a command in a container",
	Long:  `Exec a command in a container. Default to bash`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify a service")
		}
		if root.Project == "" {
			return errors.New("please specify a project with -P")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug(root.GitlabUrl)
		log.Debug(root.Project)
		log.Debug(target)

		f, err := root.Factory()
		if err != nil {
			return err
		}
		s := signpost.New(f.Project(root.Project))
		u, err := s.Target(environment)
		if err != nil {
			return err
		}
		log.Debug(u)

		user := strings.Replace(root.Project, "/", "-", -1)
		command := []string{
			"ssh",
			"-a", "-x", "-t", "-p", "2222",
			"-l", user, u.Hostname(),
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