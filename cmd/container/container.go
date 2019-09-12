package container

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/factorysh/factory-cli/cmd/root"
)

var (
	dry_run bool
	target  string
)

func init() {
	root.FlagE(execCmd.PersistentFlags())
	execCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "D", false, "DryRun")
	root.FlagE(dumpCmd.PersistentFlags())

	containerCmd.AddCommand(execCmd)
	containerCmd.AddCommand(dumpCmd)
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
		if err := root.AssertEnvironment(); err != nil {
			return err
		}
		if len(args) == 0 {
			return errors.New("you must specify a service")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := root.SSHAddress()
		if err != nil {
			return err
		}
		command := []string{
			"ssh",
			"-a", "-x", "-t", "-p", "2222",
		}
		command = append(command, root.SSHExtraArgs()...)

		if len(args) < 1 {
			args = []string{args[0], "bash"}
		}
		command = append(command, []string{address, "exec"}...)
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

var dumpCmd = &cobra.Command{
	Use:   "dump database",
	Short: "Dump a database",
	Long:  "Dump a database",
	Args: func(cmd *cobra.Command, args []string) error {
		if root.Project == "" {
			return errors.New("please specify a project with -p")
		}
		if err := root.AssertEnvironment(); err != nil {
			return err
		}
		if len(args) == 0 {
			return errors.New("you must specify a database")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := root.SignPost()
		if err != nil {
			return err
		}
		u, err := s.EnvURL(root.Environment)
		dumpUrl := fmt.Sprintf("%s/databases/%s/dump",
			u.String(),
			args[0],
		)
		l := log.WithField("dumpUurl", dumpUrl)
		req, err := http.NewRequest("GET", dumpUrl, nil)
		if err != nil {
			l.WithError(err).Error()
			return nil
		}
		req.Header.Set("Accept", "text/event-stream")
		resp, err := s.Project.Session().Do(req)
		if err != nil {
			l.WithError(err).Error()
			return nil
		}
		l = l.WithField("status", resp.Status)
		l.Debug()
		// we should use a sse reader
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err == nil {
				l.Debug(string(line))
			} else {
				l.WithError(err).Error()
				return nil
			}
		}
		return nil
	},
}
