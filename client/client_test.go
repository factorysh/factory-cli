package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAuthenticate(t *testing.T) {
	m, err := parseAuthenticate(`realm="https://gitlab.bearstech.com/jwt/auth",service="container_registry"`)
	assert.NoError(t, err)
	assert.Equal(t, "https://gitlab.bearstech.com/jwt/auth", m["realm"])
	assert.Equal(t, "container_registry", m["service"])
}
