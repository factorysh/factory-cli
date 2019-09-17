package infos

import (
	"fmt"

	"github.com/factorysh/factory-cli/cmd/root"
	"github.com/factorysh/factory-cli/signpost"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_gitlab "github.com/xanzy/go-gitlab"
)

var (
	withEnvironments bool
	withTargets      bool
)

func init() {
	infosCmd.PersistentFlags().BoolVarP(&withEnvironments, "with-environments", "e", false, "Show environments infos")
	infosCmd.PersistentFlags().BoolVarP(&withTargets, "with-targets", "T", false, "Show targets hosts")
	root.RootCmd.AddCommand(infosCmd)
}

var infosCmd = &cobra.Command{
	Use:   "infos",
	Short: "Show project's infos",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%-10v: %v\n", "project", root.Project)
		fmt.Printf("%-10v: https://%v/%v\n", "url", root.GitlabUrl, root.Project)

		if !withEnvironments {
			return nil
		}

		git, err := root.GitlabClient()
		if err != nil {
			return err
		}
		log.Debug(root.Project)

		f, err := root.Factory()
		if err != nil {
			return err
		}

		environments, _, err := git.Environments.ListEnvironments(root.Project, &_gitlab.ListEnvironmentsOptions{})
		if err != nil {
			return err
		}

		if len(environments) == 0 {
			fmt.Println("environments: none")
			return nil
		}

		s := signpost.New(f.Project(root.Project))

		log.Debug("environments: ", len(environments))

		// table headers
		fmt.Printf("\n%-15v", "environment")
		if withTargets {
			fmt.Printf("%-30v", "host")
		}
		fmt.Printf("%-30v\n", "url")

		// table sep
		fmt.Printf("--------------")
		if withTargets {
			fmt.Printf("----------------------------")
		}
		fmt.Printf("----------------------------\n")

		// one row per env
		for _, env := range environments {
			fmt.Printf("%-15v", env.Name)
			if withTargets {
				t, err := s.Target(env.Name)
				hostname := "none"
				if err == nil {
					hostname = t.Hostname()
				}
				fmt.Printf("%-30v", hostname)
			}
			if env.ExternalURL != "" {
				fmt.Printf("%-30v", env.ExternalURL)
			}
			fmt.Println("")
		}
		return nil
	},
}
