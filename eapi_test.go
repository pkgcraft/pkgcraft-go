package pkgcraft_test

import (
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

func TestEapiRange(t *testing.T) {
	// valid
	eapis, err := EapiRange("..")
	assert.Equal(t, len(eapis), len(EAPIS))
	assert.Nil(t, err)
	eapis, err = EapiRange("..=L")
	assert.Equal(t, len(eapis), len(EAPIS_OFFICIAL))
	assert.Nil(t, err)
	eapis, err = EapiRange("2..3")
	assert.Equal(t, eapis, []*Eapi{EAPIS["2"]})
	assert.Nil(t, err)
	eapis, err = EapiRange("1..")
	assert.Nil(t, err)
	assert.Equal(t, eapis[0], EAPIS["1"])
	assert.Equal(t, len(eapis), len(EAPIS)-1)
	eapis, err = EapiRange("0..0")
	assert.Nil(t, err)
	assert.Equal(t, len(eapis), 0)

	// invalid
	_, err = EapiRange("")
	assert.NotNil(t, err)
	_, err = EapiRange("1")
	assert.NotNil(t, err)
}

func TestEapiHas(t *testing.T) {
	assert.False(t, EAPIS["1"].Has("nonexistent_feature"))
	assert.True(t, EAPIS["1"].Has("slot_deps"))
}

func TestEapiString(t *testing.T) {
	for id, eapi := range EAPIS {
		assert.Equal(t, eapi.String(), id)
	}
}

func TestEapiCmp(t *testing.T) {
	assert.Equal(t, EAPIS["1"].Cmp(EAPIS["2"]), -1)
	assert.Equal(t, EAPIS["2"].Cmp(EAPIS["2"]), 0)
	assert.Equal(t, EAPIS["2"].Cmp(EAPIS["1"]), 1)
}
