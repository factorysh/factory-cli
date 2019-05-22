package cmd

import (
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
	containerCmd.PersistentFlags().StringVarP(&target, "target", "H", "localhost", "Host")
	containerCmd.PersistentFlags().StringVarP(&project, "project", "P", project_path, "Project")
	containerCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")
	rootCmd.AddCommand(containerCmd)
}

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Do something on a container",
	Long: `"Do something on a container.
factory container exec service`,

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
			"exec", "web", "bash",
		}

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
