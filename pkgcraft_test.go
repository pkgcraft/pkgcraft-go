package pkgcraft

// #cgo pkg-config: pkgcraft

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestAtom(t *testing.T) {
	var atom *Atom

	// unversioned
	atom, _ = NewAtom("cat/pkg")
	assert.Equal(t, atom.category(), "cat")
	assert.Equal(t, atom.pn(), "pkg")
	assert.Equal(t, atom.version(), "")
	assert.Equal(t, fmt.Sprintf("%s", atom), "cat/pkg")

	// versioned
	atom, _ = NewAtom("=cat/pkg-2")
	assert.Equal(t, atom.category(), "cat")
	assert.Equal(t, atom.pn(), "pkg")
	assert.Equal(t, atom.version(), "2")
	assert.Equal(t, fmt.Sprintf("%s", atom), "=cat/pkg-2")

	// slotted
	atom, _ = NewAtom("cat/pkg:1")
	assert.Equal(t, atom.slot(), "1")
	assert.Equal(t, atom.subslot(), "")
	assert.Equal(t, fmt.Sprintf("%s", atom), "cat/pkg:1")

	// subslotted
	atom, _ = NewAtom("cat/pkg:1/2")
	assert.Equal(t, atom.slot(), "1")
	assert.Equal(t, atom.subslot(), "2")
	assert.Equal(t, fmt.Sprintf("%s", atom), "cat/pkg:1/2")

	// a1 < a2
	a1, _ := NewAtom("=cat/pkg-1")
	a2, _ := NewAtom("=cat/pkg-2")
	assert.Equal(t, a1.cmp(a2), -1)

	// a1 == a2
	a1, _ = NewAtom("=cat/pkg-2")
	a2, _ = NewAtom("=cat/pkg-2")
	assert.Equal(t, a1.cmp(a2), 0)

	// a1 > a2
	a1, _ = NewAtom("=cat/pkg-2")
	a2, _ = NewAtom("=cat/pkg-1")
	assert.Equal(t, a1.cmp(a2), 1)
}

func BenchmarkNewAtom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewAtom("=cat/pkg-1-r2:3/4=[a,b,c]")
	}
}

func BenchmarkAtomSort(b *testing.B) {
	var atoms []*Atom
	for i := 100; i > 0; i-- {
		a, _ := NewAtom(fmt.Sprintf("=cat/pkg-%d", i))
		atoms = append(atoms, a)
	}
	assert.Equal(b, fmt.Sprintf("%s", atoms[0]), "=cat/pkg-100")
	for i := 0; i < b.N; i++ {
		sort.Sort(Atoms(atoms))
	}
	assert.Equal(b, fmt.Sprintf("%s", atoms[0]), "=cat/pkg-1")
}

func TestVersion(t *testing.T) {
	var version *Version

	// non-revision
	version, _ = NewVersion("1")
	assert.Equal(t, version.revision(), "0")
	assert.Equal(t, fmt.Sprintf("%s", version), "1")

	// revisioned
	version, _ = NewVersion("1-r1")
	assert.Equal(t, version.revision(), "1")
	assert.Equal(t, fmt.Sprintf("%s", version), "1-r1")

	// explicit '0' revision
	version, _ = NewVersion("1-r0")
	assert.Equal(t, version.revision(), "0")
	assert.Equal(t, fmt.Sprintf("%s", version), "1-r0")

	// v1 < v2
	v1, _ := NewVersion("1")
	v2, _ := NewVersion("2")
	assert.Equal(t, v1.cmp(v2), -1)

	// v1 == v2
	v1, _ = NewVersion("2")
	v2, _ = NewVersion("2")
	assert.Equal(t, v1.cmp(v2), 0)

	// v1 > v2
	v1, _ = NewVersion("2")
	v2, _ = NewVersion("1")
	assert.Equal(t, v1.cmp(v2), 1)
}

func BenchmarkNewVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewVersion("1.2.3_alpha4-r5")
	}
}

func BenchmarkVersionSort(b *testing.B) {
	var versions []*Version
	for i := 100; i > 0; i-- {
		v, _ := NewVersion(fmt.Sprintf("%d", i))
		versions = append(versions, v)
	}
	assert.Equal(b, fmt.Sprintf("%s", versions[0]), "100")
	for i := 0; i < b.N; i++ {
		sort.Sort(Versions(versions))
	}
	assert.Equal(b, fmt.Sprintf("%s", versions[0]), "1")
}
