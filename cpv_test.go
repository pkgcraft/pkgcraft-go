package pkgcraft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestCpv(t *testing.T) {
	var cpv *Cpv
	var ver *Version

	// valid
	cpv, _ = NewCpv("cat/pkg-1-r2")
	assert.Equal(t, cpv.Category(), "cat")
	assert.Equal(t, cpv.Package(), "pkg")
	ver, _ = NewVersion("1-r2")
	assert.Equal(t, cpv.Version(), ver)
	assert.Equal(t, cpv.Revision(), "2")
	assert.Equal(t, cpv.Key(), "cat/pkg")
	assert.Equal(t, cpv.String(), "cat/pkg-1-r2")

	cpv, _ = NewCpv("cat/pkg-0-r0")
	ver, _ = NewVersion("0-r0")
	assert.Equal(t, cpv.Version(), ver)
	assert.Equal(t, cpv.Revision(), "0")
	assert.Equal(t, cpv.String(), "cat/pkg-0-r0")

	// invalid
	_, err := NewCpv("=cat/pkg-1-r2")
	assert.NotNil(t, err)
}
