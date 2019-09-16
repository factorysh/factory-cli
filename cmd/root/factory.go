package root

import (
	"github.com/factorysh/factory-cli/factory"
)

func Factory() (*factory.Factory, error) {
	return factory.New(client, GitlabUrl, GitlabToken)
}
