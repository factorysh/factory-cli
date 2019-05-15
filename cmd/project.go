package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	_gitlab "github.com/xanzy/go-gitlab"
)

func init() {
	rootCmd.AddCommand(projectsCmd)
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Show projects",
	RunE: func(cmd *cobra.Command, args []string) error {

		git := _gitlab.NewClient(nil, os.Getenv("PRIVATE_TOKEN"))
		git.SetBaseURL(fmt.Sprintf("https://%s/api/v4", gitlab))
		projects, _, err := git.Projects.ListProjects(&_gitlab.ListProjectsOptions{})
		if err != nil {
			return err
		}

		for _, p := range projects {
			fmt.Println(p.PathWithNamespace)
		}
		return nil
	},
}
