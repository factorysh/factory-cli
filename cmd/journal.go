package cmd

import (
	"fmt"
	"os"

	"github.com/factorysh/go-longrun/longrun/sse"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_gitlab "gitlab.bearstech.com/factory/factory-cli/gitlab"
	"gitlab.bearstech.com/factory/factory-cli/journaleux"
)

var (
	lines  int
	target string
)

func init() {
	journalCmd.PersistentFlags().IntVarP(&lines, "lines", "n", -10, "Number of lines to display")
	journalCmd.PersistentFlags().StringVarP(&target, "target", "H", "localhost", "Host")
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
		j.Project(project).Logs(&journaleux.LogsOpt{
			Project: project,
			Lines:   lines,
		}, func(evt *sse.Event) error {
			fmt.Println(evt)
			return nil
		})
		return nil
	},
}
