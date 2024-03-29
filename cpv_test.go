package pkgcraft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestNewCpv(t *testing.T) {
	var cpn *Cpn
	var cpv *Cpv
	var err error
	var ver *Version
	var rev *Revision

	// valid
	cpv, err = NewCpv("cat/pkg-1-r2")
	assert.Nil(t, err)
	assert.Equal(t, cpv.Category(), "cat")
	assert.Equal(t, cpv.Package(), "pkg")
	ver, _ = NewVersion("1-r2")
	assert.Equal(t, cpv.Version(), ver)
	rev, _ = NewRevision("2")
	assert.Equal(t, cpv.Revision(), rev)
	assert.Equal(t, cpv.P(), "pkg-1")
	assert.Equal(t, cpv.Pf(), "pkg-1-r2")
	assert.Equal(t, cpv.Pr(), "r2")
	assert.Equal(t, cpv.Pv(), "1")
	assert.Equal(t, cpv.Pvr(), "1-r2")
	cpn, _ = NewCpn("cat/pkg")
	assert.Equal(t, cpv.Cpn(), cpn)
	assert.Equal(t, cpv.String(), "cat/pkg-1-r2")

	cpv, err = NewCpv("cat/pkg-0-r0")
	assert.Nil(t, err)
	ver, _ = NewVersion("0-r0")
	assert.Equal(t, cpv.Version(), ver)
	rev, _ = NewRevision("0")
	assert.Equal(t, cpv.Revision(), rev)
	assert.Equal(t, cpv.String(), "cat/pkg-0-r0")

	// invalid
	_, err = NewCpv("=cat/pkg-1-r2")
	assert.NotNil(t, err)
}

func TestCpvCmp(t *testing.T) {
	// c1 < c2
	c1, _ := NewCpv("cat/pkg-1")
	c2, _ := NewCpv("cat/pkg-2")
	assert.Equal(t, c1.Cmp(c2), -1)

	// c1 == c2
	c1, _ = NewCpv("cat/pkg-2")
	c2, _ = NewCpv("cat/pkg-2")
	assert.Equal(t, c1.Cmp(c2), 0)

	// c1 > c2
	c1, _ = NewCpv("cat/pkg-2")
	c2, _ = NewCpv("cat/pkg-1")
	assert.Equal(t, c1.Cmp(c2), 1)
}
