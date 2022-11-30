package pkgcraft_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestEapiGlobals(t *testing.T) {
	assert.True(t, len(EAPIS) > len(EAPIS_OFFICIAL))
	for id, eapi := range EAPIS_OFFICIAL {
		assert.True(t, EAPIS[id] == eapi)
	}
	assert.True(t, EAPIS[EAPI_LATEST.String()] == EAPI_LATEST)
}

func TestEapiHas(t *testing.T) {
	assert.False(t, EAPIS["1"].Has("nonexistent_feature"))
	assert.True(t, EAPIS["1"].Has("slot_deps"))
}

func TestEapiString(t *testing.T) {
	for id, eapi := range EAPIS {
		assert.Equal(t, fmt.Sprintf("%s", eapi), id)
	}
}

func TestEapiRange(t *testing.T) {
	eapis, _ := EapiRange("2-3")
	assert.Equal(t, eapis, []*Eapi{EAPIS["2"], EAPIS["3"]})
	eapis, _ = EapiRange("1-")
	assert.Equal(t, eapis[0], EAPIS["1"])
	assert.Equal(t, len(eapis), len(EAPIS)-1)
}
