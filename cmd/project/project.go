package project

import (
	"fmt"

	"github.com/factorysh/factory-cli/cmd/root"
	"github.com/factorysh/factory-cli/signpost"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_gitlab "github.com/xanzy/go-gitlab"
)

func init() {
	projectCmd.AddCommand(projectLsCmd)
	projectCmd.AddCommand(environmentsCmd)
	projectCmd.AddCommand(projectTargetCmd)
	root.RootCmd.AddCommand(projectCmd)
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Projects subcommand",
}

var projectLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		git, err := root.GitlabClient()
		if err != nil {
			return err
		}
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
			for _, p := range projects {
				fmt.Println(p.PathWithNamespace)
			}
			if r.CurrentPage < r.TotalPages {
				page++
			} else {
				break
			}
		}
		return nil
	},
}

var environmentsCmd = &cobra.Command{
	Use:   "environments",
	Short: "Show environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		git, err := root.GitlabClient()
		if err != nil {
			return err
		}
		log.Debug(root.Project)
		environments, _, err := git.Environments.ListEnvironments(root.Project, &_gitlab.ListEnvironmentsOptions{})
		if err != nil {
			return err
		}
		log.Debug("environments: ", len(environments))
		for _, env := range environments {
			if env.ExternalURL != "" {
				fmt.Printf("%v: %v\n", env.Name, env.ExternalURL)
			} else {
				fmt.Printf("%v\n", env.Name)
			}
		}
		return nil
	},
}

var projectTargetCmd = &cobra.Command{
	Use:   "target",
	Short: "Show targets of a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("[project] environment")
		}
		environment := args[0]
		log.WithField("project", root.Project).
			WithField("environment", environment).
			WithField("gitlab_url", root.GitlabUrl).
			Debug()
		f, err := root.Factory()
		if err != nil {
			return err
		}
		t, err := signpost.New(f.Project(root.Project)).Target(environment)
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	},
}
