package pkgcraft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestConfig(t *testing.T) {
	var config *Config

	// empty
	config, _ = NewConfig()
	assert.Equal(t, config.Repos().Len(), 0)
}
