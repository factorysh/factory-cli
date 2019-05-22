package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_gitlab "github.com/xanzy/go-gitlab"
	"gitlab.bearstech.com/factory/factory-cli/factory"
	__gitlab "gitlab.bearstech.com/factory/factory-cli/gitlab"
	"gitlab.bearstech.com/factory/factory-cli/signpost"
)

func init() {
	projectCmd.AddCommand(projectLsCmd)
	projectCmd.AddCommand(environmentsCmd)
	projectCmd.AddCommand(projectTargetCmd)
	rootCmd.AddCommand(projectCmd)
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Projects subcommand",
}

var projectLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects",
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

var environmentsCmd = &cobra.Command{
	Use:   "environments",
	Short: "Show environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		git := _gitlab.NewClient(nil, os.Getenv("PRIVATE_TOKEN"))
		gitlab, err := guessGitlab()
		if err != nil {
			return nil
		}
		git.SetBaseURL(fmt.Sprintf("https://%s/api/v4", gitlab))
		var project string
		if len(args) > 0 {
			project = args[0]
		} else {
			_, project, err = __gitlab.GitRemote()
			if err != nil {
				return err
			}
		}
		log.Debug(project)
		environments, _, err := git.Environments.ListEnvironments(project, &_gitlab.ListEnvironmentsOptions{})
		if err != nil {
			return err
		}
		for _, env := range environments {
			fmt.Printf("%v: %v\n", env.Name, env.ExternalURL)
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
		gitlab, err := guessGitlab()
		if err != nil {
			return nil
		}
		var (
			project, environment string
		)
		if len(args) == 1 {
			_, project, err = __gitlab.GitRemote()
			if err != nil {
				return err
			}
			environment = args[0]
		} else {
			project = args[0]
			environment = args[1]
		}
		log.WithField("project", project).
			WithField("environment", environment).
			WithField("gitlab", gitlab).
			Debug()
		f, err := factory.New(gitlab, os.Getenv("PRIVATE_TOKEN"))
		if err != nil {
			return err
		}
		t, err := signpost.New(f.Project(project)).Target(environment)
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	},
}
