package volume

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
	environment string
	dry_run     bool
)

func init() {
	sftpCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "Environment")
	sftpCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")
	volumeCmd.AddCommand(sftpCmd)

	root.RootCmd.AddCommand(volumeCmd)
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Do something on a volume",
	Long:  `Do something on a volume.`,
}

var sftpCmd = &cobra.Command{
	Use:   "sftp",
	Short: "sftp to project's volumes",
	Long:  `sftp to project's volumes`,
	Args: func(cmd *cobra.Command, args []string) error {
		if root.Project == "" {
			return errors.New("please specify a project with -P")
		}
		return cobra.NoArgs(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug(root.GitlabUrl)
		log.Debug(root.Project)

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
			"sftp",
			"-P", "2222",
			user + "@" + u.Hostname(),
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
