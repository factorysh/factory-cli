package root

import (
	"os"

	"gitlab.bearstech.com/factory/factory-cli/factory"
)

func Factory() (*factory.Factory, error) {
	return factory.New(client, GitlabUrl, os.Getenv("PRIVATE_TOKEN"))
}
