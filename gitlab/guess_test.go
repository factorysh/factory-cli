package gitlab

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	_, project, err := GitRemote()
	assert.NoError(t, err)
	fmt.Printf("%#v", project)
	assert.True(t, strings.HasSuffix(project, "factory-cli"))
}
