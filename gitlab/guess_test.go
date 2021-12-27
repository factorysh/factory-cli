package gitlab

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	// authenticated url
	server, project := httpUrl("https://user:pass@github.com/factorysh/factory-cli.git")
	assert.True(t, strings.HasSuffix(project, "factorysh/factory-cli"))
	assert.True(t, strings.HasSuffix(server, "github.com"))

	// unauthenticated url (github actions)
	server, project = httpUrl("https://github.com/factorysh/factory-cli")
	assert.True(t, strings.HasSuffix(project, "factorysh/factory-cli"))
	assert.True(t, strings.HasSuffix(server, "github.com"))

	// ssh url (github actions)
	server, project = gitUrl("git@github.com:factorysh/factory-cli.git")
	assert.True(t, strings.HasSuffix(project, "factorysh/factory-cli"))
	assert.True(t, strings.HasSuffix(server, "github.com"))

	_, project, err := GitRemote()
	assert.NoError(t, err)
	fmt.Printf("%#v", project)
	assert.True(t, strings.HasSuffix(project, "factorysh/factory-cli"))
}
