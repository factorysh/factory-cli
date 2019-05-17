package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_gitlab "gitlab.bearstech.com/factory/factory-cli/gitlab"
	"gitlab.bearstech.com/factory/factory-cli/journaleux"
)

var (
	lines  int
	target string
	format string
)

func init() {
	journalCmd.PersistentFlags().IntVarP(&lines, "lines", "n", -10, "Number of lines to display")
	journalCmd.PersistentFlags().StringVarP(&target, "target", "H", "localhost", "Host")
	journalCmd.PersistentFlags().StringVar(&format, "format", "bare", "Output format : bare|json")
	rootCmd.AddCommand(journalCmd)
}

var journalCmd = &cobra.Command{
	Use:   "journal",
	Short: "Show journal",

	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
		var (
			project string
			err     error
		)
		if len(args) > 0 {
			project = args[0]
		} else {
			_, project, err = _gitlab.GitRemote()
			if err != nil {
				return err
			}
		}
		j, err := journaleux.New(target, os.Getenv("PRIVATE_TOKEN"))
		if err != nil {
			return err
		}
		h, err := j.Project(project).Hello()
		if err != nil {
			return err
		}
		log.Debug(h)
		cpt := 0
		j.Project(project).Logs(&journaleux.LogsOpt{
			Project: project,
			Lines:   lines,
		}, func(evt *journaleux.Event, zerr error) error {
			switch format {
			case "bare":
				fmt.Println(evt.Message)
			case "json":
				j, err := json.Marshal(evt)
				if err != nil {
					return err
				}
				fmt.Println(string(j))
			default:
				fmt.Println(evt)
			}
			cpt++
			return nil
		})
		log.Debug("Lines:", cpt)
		return nil
	},
}
