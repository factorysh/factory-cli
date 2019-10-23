package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/factorysh/factory-cli/cmd/root"
	"github.com/factorysh/go-longrun/longrun/sse"
)

var (
	dry_run  bool
	download bool
	target   string
)

func init() {
	root.FlagE(execCmd.PersistentFlags())
	execCmd.PersistentFlags().BoolVarP(
		&dry_run, "dry-run", "D", false, "Only print ssh command")
	root.FlagE(dumpCmd.PersistentFlags())
	dumpCmd.PersistentFlags().BoolVarP(
		&download, "no-download", "", true, "Do not download the file locally")

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
	Use:   "exec <service> [command]",
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
	Use:   "dump <database>",
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
		l := log.WithField("dumpUrl", dumpUrl)
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
		filename := ""
		defer resp.Body.Close()
		err = sse.Reader(resp.Body, func(evt *sse.Event) error {
			var event map[string]interface{}
			err := json.Unmarshal([]byte(evt.Data), &event)
			if err != nil {
				l.WithError(err).Error()
				return err
			}
			// filename is stored in the first debug task
			// means that we wait for the first message like:
			// {"result": {"msg": ""}}
			res := event["result"]
			if res != nil {
				value := res.(map[string]interface{})["msg"]
				if value != nil && filename == "" {
					filename = value.(string)
					l = l.WithField("filename", filename)
					l.Debug("Got filename")
				}
			}
			return nil
		})

		if filename != "" && download {
			// sftp the file
			fmt.Println("Fetching", filename)
			_url, err := root.SSHAddress()
			if err != nil {
				log.WithError(err).Error()
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
			c := exec.Command("sftp")
			c.Args = command
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			// create a pipe for the sftp "get <filename>" command
			stdin, err := c.StdinPipe()
			if err != nil {
				log.WithError(err).Error()
				return err
			}
			if err = c.Start(); err != nil {
				log.WithError(err).Error()
				return err
			}
			// write the command to stdin
			io.WriteString(stdin, "get /volumes/snapshots/"+filename+"\n")
			stdin.Close()
			c.Wait()
		}
		return nil
	},
}
