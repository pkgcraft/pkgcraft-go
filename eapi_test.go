package pkgcraft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestEapiGlobals(t *testing.T) {
	assert.True(t, len(EAPIS) > len(EAPIS_OFFICIAL))
	for id, eapi := range EAPIS_OFFICIAL {
		assert.True(t, EAPIS[id] == eapi)
	}
	assert.True(t, EAPIS[EAPI_LATEST_OFFICIAL.String()] == EAPI_LATEST_OFFICIAL)
}

func TestEapiRange(t *testing.T) {
	// valid
	eapis, err := EapiRange("..")
	assert.Equal(t, len(eapis), len(EAPIS))
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
	for _, s := range []string{"", "1", "..9999"} {
		_, err = EapiRange(s)
		assert.NotNil(t, err)
	}
}

func TestEapiHas(t *testing.T) {
	assert.False(t, EAPIS["1"].Has("nonexistent_feature"))
	assert.True(t, EAPIS["1"].Has("slot_deps"))
}

func TestEapiDepKeys(t *testing.T) {
	assert.True(t, slices.Contains(EAPIS["0"].DepKeys(), "DEPEND"))
	assert.False(t, slices.Contains(EAPIS["0"].DepKeys(), "BDEPEND"))
	assert.True(t, slices.Contains(EAPI_LATEST.DepKeys(), "BDEPEND"))
}

func TestEapiMetadataKeys(t *testing.T) {
	assert.True(t, slices.Contains(EAPIS["0"].MetadataKeys(), "SLOT"))
	assert.False(t, slices.Contains(EAPIS["0"].MetadataKeys(), "BDEPEND"))
	assert.True(t, slices.Contains(EAPI_LATEST.MetadataKeys(), "BDEPEND"))
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
