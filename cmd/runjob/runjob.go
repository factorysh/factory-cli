package runjob

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/factorysh/factory-cli/cmd/root"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const command_label = "sh.factory.cronjob.command"

var (
	dry_run bool
)

func init() {
	runjobCmd.PersistentFlags().BoolVarP(
		&dry_run, "dry-run", "D", false, "Only print command")
	root.RootCmd.AddCommand(runjobCmd)
}

var runjobCmd = &cobra.Command{
	Use:   "runjob <service>",
	Short: "Run a job from your local project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify a job name")
		} else if len(args) > 1 {
			return errors.New("you must only specify a job name")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		compose, ok := os.LookupEnv("COMPOSE")
		if !ok {
			return errors.New("No COMPOSE env var found")
		}

		file, err := os.Open("docker-compose.yml")

		if err != nil {
			log.Fatalf("No docker-compose.yml found")
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		job := args[0]
		command := ""
		dry_run_command := ""
		in_job := false
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.Trim(line, " ")
			if strings.HasPrefix(line, job+":") {
				in_job = true
				log.Debug("Found job: " + job)
			} else if strings.Index(line, command_label) >= 0 && in_job {
				command = strings.Split(line, command_label+":")[1]
				command = strings.Trim(command, " ")
				dry_run_command = command
				if strings.HasPrefix(command, "\"") {
					command = strings.Trim(command, "\"")
				} else if strings.HasPrefix(command, "'") {
					command = strings.Trim(command, "'")
				} else {
					dry_run_command = "\"" + command + "\""
				}
				log.Debug("Found command: " + command)
				break
			}
		}

		file.Close()

		if !in_job {
			return errors.New("Service " + job + " not found in your docker-compose.yml")
		}
		if command == "" {
			return errors.New("No " + command_label + "found in service " + job)
		}

		log.Debug("Command: " + command)
		if command != "" {
			compose_command := strings.Split(compose, " ")
			dry_run_compose_command := append(
				compose_command,
				"run", "--rm", job, "bash", "-c", dry_run_command,
			)
			compose_command = append(
				compose_command,
				"run", "--rm", job, "bash", "-c", command,
			)
			if dry_run {
				fmt.Println(strings.Join(dry_run_compose_command, " "))
			} else {
				c := exec.Command(compose_command[0])
				c.Args = compose_command
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				err := c.Run()
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}
