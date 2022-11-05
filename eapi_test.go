package pkgcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEapiGlobals(t *testing.T) {
	assert.True(t, len(EAPIS) > len(EAPIS_OFFICIAL))
	for id, eapi := range EAPIS_OFFICIAL {
		assert.True(t, EAPIS[id] == eapi)
	}
	assert.True(t, EAPIS[EAPI_LATEST.String()] == EAPI_LATEST)
}

func TestEapiHas(t *testing.T) {
	eapi := EAPIS["1"]
	assert.False(t, eapi.has("nonexistent_feature"))
	assert.True(t, eapi.has("slot_deps"))
}
