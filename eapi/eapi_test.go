package eapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEapiGlobals(t *testing.T) {
	assert.True(t, len(Eapis) > len(EapisOfficial))
	for id, eapi := range EapisOfficial {
		assert.True(t, Eapis[id] == eapi)
	}
	assert.True(t, Eapis[EapiLatest.String()] == EapiLatest)
}

func TestEapiHas(t *testing.T) {
	eapi := Eapis["1"]
	assert.False(t, eapi.has("nonexistent_feature"))
	assert.True(t, eapi.has("slot_deps"))
}
