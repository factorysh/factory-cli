package root

import (
	"os"

	"github.com/factorysh/factory-cli/factory"
)

func Factory() (*factory.Factory, error) {
	return factory.New(client, GitlabUrl, os.Getenv("PRIVATE_TOKEN"))
}
