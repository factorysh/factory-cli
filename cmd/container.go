package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dry_run bool
)

func init() {
	execCmd.PersistentFlags().StringVarP(&target, "target", "H", "localhost", "Host")
	execCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")

	containerCmd.AddCommand(execCmd)
	rootCmd.AddCommand(containerCmd)
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
		if project == "" {
			return errors.New("please specify a project with -P")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug(project)
		log.Debug(target)

		user := strings.Replace(project, "/", "-", -1)
		command := []string{
			"ssh",
			"-a", "-x", "-t", "-p", "2222",
			"-l", user, target,
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
