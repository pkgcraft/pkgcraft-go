package pkgcraft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestNewCpn(t *testing.T) {
	var cpn *Cpn
	var err error

	// valid
	cpn, err = NewCpn("cat/pkg")
	assert.Nil(t, err)
	assert.Equal(t, cpn.Category(), "cat")
	assert.Equal(t, cpn.Package(), "pkg")

	cpn, err = NewCpn("cat/pkg-1-a")
	assert.Nil(t, err)
	assert.Equal(t, cpn.Category(), "cat")
	assert.Equal(t, cpn.Package(), "pkg-1-a")

	// invalid
	_, err = NewCpn("cat/pkg-1")
	assert.NotNil(t, err)
	_, err = NewCpn("=cat/pkg-1-r2")
	assert.NotNil(t, err)
}

func TestCpnCmp(t *testing.T) {
	// c1 < c2
	c1, _ := NewCpn("a/b")
	c2, _ := NewCpn("a/c")
	assert.Equal(t, c1.Cmp(c2), -1)

	// c1 == c2
	c1, _ = NewCpn("a/b")
	c2, _ = NewCpn("a/b")
	assert.Equal(t, c1.Cmp(c2), 0)

	// c1 > c2
	c1, _ = NewCpn("b/a")
	c2, _ = NewCpn("a/b")
	assert.Equal(t, c1.Cmp(c2), 1)
}
