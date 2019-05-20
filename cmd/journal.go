package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_gitlab "gitlab.bearstech.com/factory/factory-cli/gitlab"
	"gitlab.bearstech.com/factory/factory-cli/journaleux"
)

var (
	lines     int
	target    string
	format    string
	timestamp bool
	follow    bool
	re        string
)

func init() {
	journalCmd.PersistentFlags().IntVarP(&lines, "lines", "n", -10, "Number of lines to display")
	journalCmd.PersistentFlags().StringVarP(&target, "target", "H", "localhost", "Host")
	journalCmd.PersistentFlags().StringVar(&format, "format", "bare", "Output format : bare|json|jsonpretty")
	journalCmd.PersistentFlags().BoolVarP(&timestamp, "timestamp", "t", false, "Show timestamps")
	journalCmd.PersistentFlags().BoolVarP(&follow, "follow", "f", false, "Follow")
	journalCmd.PersistentFlags().StringVarP(&re, "regexp", "r", "", "Regular expression filter")
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
		if re != "" {
			_, err = regexp.Compile(re)
			if err != nil {
				return err
			}
		}
		j.Project(project).Logs(&journaleux.LogsOpt{
			Project: project,
			Lines:   lines,
			Follow:  follow,
			Regexp:  re,
		}, func(evt *journaleux.Event, zerr error) error {
			switch format {
			case "bare":
				if timestamp {
					t := time.Unix(int64(evt.Realtime)/1000000, (int64(evt.Realtime)%1000000)*1000)
					fmt.Print(t.Format(time.RFC3339), " ")
				}
				fmt.Println(evt.Message)
			case "json":
				j, err := json.Marshal(evt)
				if err != nil {
					return err
				}
				fmt.Println(string(j))
			case "jsonpretty":
				j, err := json.MarshalIndent(evt, "", "  ")
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
