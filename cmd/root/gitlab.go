package root

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

func GitlabClient() (*gitlab.Client, error) {
	git := gitlab.NewClient(client, os.Getenv("PRIVATE_TOKEN"))
	git.SetBaseURL(fmt.Sprintf("https://%s/api/v4", GitlabUrl))
	return git, nil
}

func Environments() ([]string, error) {
	git, err := GitlabClient()
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}
	l := log.WithField("project", Project)
	environments, _, err := git.Environments.ListEnvironments(Project, &gitlab.ListEnvironmentsOptions{})
	if err != nil {
		l.WithError(err).Error()
		return nil, err
	}
	envs := make([]string, len(environments))
	for i, env := range environments {
		envs[i] = env.Name
	}
	l.WithField("environments", envs).Debug()
	return envs, nil
}

func AssertEnvironment() error {
	if Environment != "" {
		return nil
	}
	envs, err := Environments()
	if err != nil {
		return err
	}
	if len(envs) == 0 {
		return fmt.Errorf("You are doomed, the project %s has no environment", Project)
	}
	return fmt.Errorf("Select an environment for the project %s with -e option: %s",
		Project, strings.Join(envs, "|"))
}
