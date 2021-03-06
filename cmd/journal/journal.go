package journal

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/factorysh/factory-cli/cmd/root"
	"github.com/factorysh/factory-cli/journaleux"
	"github.com/factorysh/factory-cli/signpost"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	lines     int
	format    string
	timestamp bool
	follow    bool
	re        string
	fields    map[string]string
)

func init() {
	journalCmd.PersistentFlags().IntVarP(&lines, "lines", "n", -10, "Number of lines to display")
	journalCmd.PersistentFlags().StringVar(&format, "format", "bare", "Output format : bare|json|jsonpretty")
	journalCmd.PersistentFlags().BoolVarP(&timestamp, "timestamp", "T", false, "Show timestamps")
	journalCmd.PersistentFlags().BoolVarP(&follow, "follow", "f", false, "Follow")
	journalCmd.PersistentFlags().StringVarP(&re, "regexp", "r", "", "Regular expression filter")
	root.FlagE(journalCmd.PersistentFlags())
	root.RootCmd.AddCommand(journalCmd)
}

var journalCmd = &cobra.Command{
	Use:     "journal [key=value …]",
	Short:   "Show journal",
	Long:    `Show journal of a project.`,
	Example: `factory journal -n 100 COM_DOCKER_COMPOSE_SERVICE=cron`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := root.AssertEnvironment()
		if err != nil {
			return err
		}
		fields, err = guessArgs(args)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := root.Factory()
		if err != nil {
			return err
		}

		t, err := signpost.New(f.Project(root.Project)).Target(root.Environment)
		if err != nil {
			return err
		}
		t.Host = t.Hostname()
		log.Debug("target: ", t.Host)
		j := journaleux.New(f.Project(root.Project), t)
		cpt := 0
		if re != "" {
			_, err = regexp.Compile(re)
			if err != nil {
				return err
			}
		}
		j.Logs(&journaleux.LogsOpt{
			Project: root.Project,
			Lines:   lines,
			Follow:  follow,
			Regexp:  re,
			Fields:  fields,
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

func guessArgs(args []string) (fields map[string]string, err error) {
	if len(args) == 0 {
		return nil, nil
	}
	fields = make(map[string]string)
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("Bad key=value format: %v", arg)
		}
		fields[kv[0]] = kv[1]
	}
	return fields, nil
}
