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

func init() {
	_, project_path, _ := _gitlab.GitRemote()
	sftpCmd.PersistentFlags().StringVarP(&target, "target", "H", "localhost", "Host")
	volumeCmd.AddCommand(sftpCmd)

	volumeCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")
	volumeCmd.PersistentFlags().StringVarP(&project, "project", "P", project_path, "Project")
	rootCmd.AddCommand(volumeCmd)
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Do something on a volume",
	Long:  `Do something on a volume.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
		return errors.New("you must use a subcommand: sftp")
	},
}

var sftpCmd = &cobra.Command{
	Use:   "sftp",
	Short: "sftp to project's volumes",
	Long:  `sftp to project's volumes`,
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
			"sftp",
			"-P", "2222",
			user + "@" + target,
		}
		command = append(command, args...)

		log.Debug(command)
		if dry_run {
			fmt.Printf("%s\n", strings.Join(command, " "))
		} else {
			c := exec.Command("sftp")
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
