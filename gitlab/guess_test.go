package gitlab

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	_, project, err := GitRemote()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(project, "factory-cli"))
}
