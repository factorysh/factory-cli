package volume

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/factorysh/factory-cli/cmd/root"
)

var (
	environment string
	dry_run     bool
)

func init() {
	root.FlagE(sftpCmd.PersistentFlags())
	sftpCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")
	root.FlagE(volumeCmd.PersistentFlags())
	volumeCmd.AddCommand(sftpCmd)
	volumeCmd.AddCommand(urlCmd)

	root.RootCmd.AddCommand(volumeCmd)
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Do something on a volume",
	Long:  `Do something on a volume.`,
}

var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "url of the volume",
	Args: func(cmd *cobra.Command, args []string) error {
		if root.Project == "" {
			return errors.New("please specify a project with -p")
		}
		return root.AssertEnvironment()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		u, err := root.SSHAddress()
		if err != nil {
			return err
		}
		fmt.Printf("sftp://%s:2222/\n", u)
		return nil
	},
}

var sftpCmd = &cobra.Command{
	Use:   "sftp",
	Short: "sftp to project's volumes",
	Long:  `sftp to project's volumes`,
	Args: func(cmd *cobra.Command, args []string) error {
		if root.Project == "" {
			return errors.New("please specify a project with -p")
		}
		if err := root.AssertEnvironment(); err != nil {
			return err
		}
		return cobra.NoArgs(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_url, err := root.SSHAddress()
		if err != nil {
			return err
		}

		command := []string{
			"sftp",
			"-P", "2222",
		}

		command = append(command, root.SSHExtraArgs()...)

		command = append(command, _url)

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
