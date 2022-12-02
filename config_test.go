package pkgcraft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestConfig(t *testing.T) {
	// empty
	config, _ := NewConfig()
	defer config.Close()
	assert.Equal(t, len(config.Repos), 0)

	// verify repo maps are empty
	_, exists := config.Repos["repo"]
	assert.False(t, exists)
	_, exists = config.ReposEbuild["repo"]
	assert.False(t, exists)
	_, exists = config.ReposFake["repo"]
	assert.False(t, exists)
}
