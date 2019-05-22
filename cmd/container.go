package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_gitlab "gitlab.bearstech.com/factory/factory-cli/gitlab"
)

var (
	project string
	dry_run bool
)

func init() {
	_, project_path, _ := _gitlab.GitRemote()
	execCmd.PersistentFlags().StringVarP(&target, "target", "H", "localhost", "Host")
	containerCmd.AddCommand(execCmd)

	containerCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")
	containerCmd.PersistentFlags().StringVarP(&project, "project", "P", project_path, "Project")
	rootCmd.AddCommand(containerCmd)
}

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Do something on a container",
	Long:  `Do something on a container.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
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
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
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
