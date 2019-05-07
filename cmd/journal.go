package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.bearstech.com/factory/factory-cli/journaleux"
)

func init() {
	rootCmd.AddCommand(journalCmd)
}

var journalCmd = &cobra.Command{
	Use:   "journal",
	Short: "Show journal",

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(args)
		log.SetLevel(log.DebugLevel)
		j, err := journaleux.New(args[0], os.Getenv("PRIVATE_TOKEN"))
		if err != nil {
			return err
		}
		h, err := j.Project(args[1]).Hello()
		if err != nil {
			return err
		}
		fmt.Println(h)
		return nil
	},
}
