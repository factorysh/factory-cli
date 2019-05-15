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
		gitlab, err := guessGitlab()
		if err != nil {
			return nil
		}
		git.SetBaseURL(fmt.Sprintf("https://%s/api/v4", gitlab))
		page := 0
		for {
			opts := &_gitlab.ListProjectsOptions{
				OrderBy: _gitlab.String("name"),
				Sort:    _gitlab.String("asc"),
			}
			opts.Page = page
			projects, r, err := git.Projects.ListProjects(opts)
			if err != nil {
				return err
			}
			if r.CurrentPage < r.TotalPages {
				page++
			} else {
				break
			}
			for _, p := range projects {
				fmt.Println(p.PathWithNamespace)
			}
		}
		return nil
	},
}
