package root

import (
	"fmt"
	"os"

	gitlab "github.com/xanzy/go-gitlab"
)

func GitlabClient() (*gitlab.Client, error) {
	git := gitlab.NewClient(client, os.Getenv("PRIVATE_TOKEN"))
	git.SetBaseURL(fmt.Sprintf("https://%s/api/v4", GitlabUrl))
	return git, nil
}
