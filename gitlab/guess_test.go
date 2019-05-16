package gitlab

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	server, project := httpUrl("https://user:pass@gitlab.com/factory/factory-cli.git")
	assert.True(t, strings.HasSuffix(project, "factory/factory-cli"))
	assert.True(t, strings.HasSuffix(server, "gitlab.com"))

	server, project = gitUrl("gitlab@gitlab.com:factory/factory-cli.git")
	assert.True(t, strings.HasSuffix(project, "factory/factory-cli"))
	assert.True(t, strings.HasSuffix(server, "gitlab.com"))

	_, project, err := GitRemote()
	assert.NoError(t, err)
	fmt.Printf("%#v", project)
	assert.True(t, strings.HasSuffix(project, "factory/factory-cli"))
}
