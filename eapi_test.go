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
	eapis, err = EapiRange("7..8")
	assert.Equal(t, eapis, []*Eapi{EAPIS["7"]})
	assert.Nil(t, err)
	eapis, err = EapiRange("7..")
	assert.Nil(t, err)
	assert.Equal(t, eapis[0], EAPIS["7"])
	eapis, err = EapiRange("8..8")
	assert.Nil(t, err)
	assert.Equal(t, len(eapis), 0)

	// invalid
	for _, s := range []string{"", "1", "..9999"} {
		_, err = EapiRange(s)
		assert.NotNil(t, err)
	}
}

func TestEapiHas(t *testing.T) {
	assert.True(t, EAPI_LATEST_OFFICIAL.Has("UsevTwoArgs"))
	assert.False(t, EAPI_LATEST_OFFICIAL.Has("RepoIds"))
	assert.False(t, EAPI_LATEST_OFFICIAL.Has("nonexistent"))
}

func TestEapiDepKeys(t *testing.T) {
	assert.True(t, slices.Contains(EAPI_LATEST.DepKeys(), "BDEPEND"))
	assert.False(t, slices.Contains(EAPI_LATEST.DepKeys(), "NONEXISTENT"))
}

func TestEapiMetadataKeys(t *testing.T) {
	assert.True(t, slices.Contains(EAPI_LATEST.MetadataKeys(), "SLOT"))
	assert.False(t, slices.Contains(EAPI_LATEST.MetadataKeys(), "NONEXISTENT"))
}

func TestEapiString(t *testing.T) {
	for id, eapi := range EAPIS {
		assert.Equal(t, eapi.String(), id)
	}
}

func TestEapiCmp(t *testing.T) {
	assert.Equal(t, EAPIS["7"].Cmp(EAPIS["8"]), -1)
	assert.Equal(t, EAPIS["8"].Cmp(EAPIS["8"]), 0)
	assert.Equal(t, EAPIS["8"].Cmp(EAPIS["7"]), 1)
}
