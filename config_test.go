package pkgcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var config *Config

	// empty
	config, _ = NewConfig()
	assert.Equal(t, config.repos().len(), 0)
}
