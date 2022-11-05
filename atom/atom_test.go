package atom

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtom(t *testing.T) {
	var atom, c1, c2 *Atom
	var ver *Version

	// unversioned
	atom, _ = NewAtom("cat/pkg")
	assert.Equal(t, atom.category(), "cat")
	assert.Equal(t, atom.pn(), "pkg")
	assert.Nil(t, atom.version())
	assert.Equal(t, atom.revision(), "")
	assert.Equal(t, atom.blocker(), BlockerNone)
	assert.Equal(t, atom.slot(), "")
	assert.Equal(t, atom.subslot(), "")
	assert.Equal(t, atom.slot_op(), SlotOpNone)
	assert.Equal(t, atom.use_deps(), []string{})
	assert.Equal(t, atom.repo(), "")
	assert.Equal(t, atom.key(), "cat/pkg")
	assert.Equal(t, atom.cpv(), "cat/pkg")
	assert.Equal(t, fmt.Sprintf("%s", atom), "cat/pkg")

	// versioned
	atom, _ = NewAtom("=cat/pkg-1-r2")
	assert.Equal(t, atom.category(), "cat")
	assert.Equal(t, atom.pn(), "pkg")
	ver, _ = NewVersionWithOp("=1-r2")
	assert.Equal(t, atom.version(), ver)
	assert.Equal(t, atom.revision(), "2")
	assert.Equal(t, atom.key(), "cat/pkg")
	assert.Equal(t, atom.cpv(), "cat/pkg-1-r2")
	assert.Equal(t, fmt.Sprintf("%s", atom), "=cat/pkg-1-r2")

	// blocker
	atom, _ = NewAtom("!cat/pkg")
	assert.Equal(t, atom.blocker(), BlockerWeak)
	assert.Equal(t, fmt.Sprintf("%s", atom), "!cat/pkg")

	// subslotted
	atom, _ = NewAtom("cat/pkg:1/2")
	assert.Equal(t, atom.slot(), "1")
	assert.Equal(t, atom.subslot(), "2")
	assert.Equal(t, fmt.Sprintf("%s", atom), "cat/pkg:1/2")

	// slot operator
	atom, _ = NewAtom("cat/pkg:0=")
	assert.Equal(t, atom.slot(), "0")
	assert.Equal(t, atom.slot_op(), SlotOpEqual)
	assert.Equal(t, fmt.Sprintf("%s", atom), "cat/pkg:0=")

	// repo
	atom, _ = NewAtom("cat/pkg::repo")
	assert.Equal(t, atom.repo(), "repo")
	assert.Equal(t, fmt.Sprintf("%s", atom), "cat/pkg::repo")

	// repo dep invalid on official EAPIs
	atom, _ = NewAtomWithEapi("cat/pkg::repo", "8")
	assert.Nil(t, atom)

	// unknown EAPI
	atom, _ = NewAtomWithEapi("cat/pkg", "unknown")
	assert.Nil(t, atom)

	// all fields
	atom, _ = NewAtom("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.Equal(t, atom.category(), "cat")
	assert.Equal(t, atom.pn(), "pkg")
	ver, _ = NewVersionWithOp("=1-r2")
	assert.Equal(t, atom.version(), ver)
	assert.Equal(t, atom.revision(), "2")
	assert.Equal(t, atom.blocker(), BlockerStrong)
	assert.Equal(t, atom.slot(), "3")
	assert.Equal(t, atom.subslot(), "4")
	assert.Equal(t, atom.slot_op(), SlotOpEqual)
	assert.Equal(t, atom.use_deps(), []string{"a", "b", "c"})
	assert.Equal(t, atom.repo(), "repo")
	assert.Equal(t, atom.key(), "cat/pkg")
	assert.Equal(t, atom.cpv(), "cat/pkg-1-r2")
	assert.Equal(t, fmt.Sprintf("%s", atom), "!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")

	// verify cached atoms reuse objects
	c1, _ = NewCachedAtom("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.Equal(t, atom.cmp(c1), 0)
	assert.True(t, atom != c1)
	c2, _ = NewCachedAtom("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.True(t, c1 == c2)
	c1, _ = NewCachedAtomWithEapi("!!=a/b-1-r2:3/4=[a,b,c]", "8")
	c2, _ = NewCachedAtomWithEapi("!!=a/b-1-r2:3/4=[a,b,c]", "8")
	assert.True(t, c1 == c2)

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

	// hashing equal values
	a1, _ = NewAtom("=cat/pkg-1.0.2")
	a2, _ = NewAtom("=cat/pkg-1.0.2-r0")
	a3, _ := NewAtom("=cat/pkg-1.000.2")
	a4, _ := NewAtom("=cat/pkg-1.00.2-r0")
	m := make(map[uint64]bool)
	m[a1.hash()] = true
	m[a2.hash()] = true
	m[a3.hash()] = true
	m[a4.hash()] = true
	assert.Equal(t, len(m), 1)

	// hashing unequal values
	a1, _ = NewAtom("=cat/pkg-0.1")
	a2, _ = NewAtom("=cat/pkg-0.01")
	a3, _ = NewAtom("=cat/pkg-0.001")
	m = make(map[uint64]bool)
	m[a1.hash()] = true
	m[a2.hash()] = true
	m[a3.hash()] = true
	assert.Equal(t, len(m), 3)
}

// test sending Atoms over a channel
func TestAtomChannels(t *testing.T) {
	var atom *Atom

	atom_strs := make(chan string)
	atoms := make(chan *Atom)

	go func() {
		for {
			s := <-atom_strs
			atom, _ = NewAtom(s)
			atoms <- atom
		}
	}()

	var s string
	for i := 0; i < 1000; i++ {
		s = fmt.Sprintf("=cat/pkg-%d-r2:3/4=[a,b,c]", i)
		atom_strs <- s
		atom = <-atoms
		assert.Equal(t, fmt.Sprintf("%s", atom), s)
	}
}

func BenchmarkNewAtom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewAtom(fmt.Sprintf("=cat/pkg-%d-r2:3/4=[a,b,c]", i))
	}
}

func BenchmarkNewAtomStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewAtom("=cat/pkg-1-r2:3/4=[a,b,c]")
	}
}

func BenchmarkNewCachedAtomStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCachedAtom("=cat/pkg-1-r2:3/4=[a,b,c]")
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

func TestCpv(t *testing.T) {
	var cpv *Cpv
	var ver *Version

	// valid
	cpv, _ = NewCpv("cat/pkg-1-r2")
	assert.Equal(t, cpv.category(), "cat")
	assert.Equal(t, cpv.pn(), "pkg")
	ver, _ = NewVersion("1-r2")
	assert.Equal(t, cpv.version(), ver)
	assert.Equal(t, cpv.revision(), "2")
	assert.Equal(t, cpv.key(), "cat/pkg")
	assert.Equal(t, cpv.cpv(), "cat/pkg-1-r2")
	assert.Equal(t, fmt.Sprintf("%s", cpv), "cat/pkg-1-r2")

	cpv, _ = NewCpv("cat/pkg-0-r0")
	ver, _ = NewVersion("0-r0")
	assert.Equal(t, cpv.version(), ver)
	assert.Equal(t, cpv.revision(), "0")
	assert.Equal(t, fmt.Sprintf("%s", cpv), "cat/pkg-0-r0")

	// invalid
	_, err := NewCpv("=cat/pkg-1-r2")
	assert.NotNil(t, err)
}
