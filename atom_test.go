package pkgcraft_test

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
	. "github.com/pkgcraft/pkgcraft-go/internal"
)

func TestAtom(t *testing.T) {
	var atom, c1, c2 *Atom
	var err error
	var ver *Version

	// unversioned
	atom, err = NewAtom("cat/pkg")
	assert.Nil(t, err)
	assert.Equal(t, atom.Category(), "cat")
	assert.Equal(t, atom.Package(), "pkg")
	assert.Equal(t, atom.Version(), &Version{})
	assert.Equal(t, atom.Revision(), "")
	assert.Equal(t, atom.Blocker(), BlockerNone)
	assert.Equal(t, atom.Slot(), "")
	assert.Equal(t, atom.Subslot(), "")
	assert.Equal(t, atom.SlotOp(), SlotOpNone)
	assert.Equal(t, len(atom.Use()), 0)
	assert.Equal(t, atom.Repo(), "")
	assert.Equal(t, atom.Cpn(), "cat/pkg")
	assert.Equal(t, atom.CPV(), "cat/pkg")
	assert.Equal(t, atom.String(), "cat/pkg")

	// versioned
	atom, err = NewAtom("=cat/pkg-1-r2")
	assert.Nil(t, err)
	assert.Equal(t, atom.Category(), "cat")
	assert.Equal(t, atom.Package(), "pkg")
	ver, _ = NewVersionWithOp("=1-r2")
	assert.Equal(t, atom.Version(), ver)
	assert.Equal(t, atom.Revision(), "2")
	assert.Equal(t, atom.Cpn(), "cat/pkg")
	assert.Equal(t, atom.CPV(), "cat/pkg-1-r2")
	assert.Equal(t, atom.String(), "=cat/pkg-1-r2")

	// blocker
	atom, err = NewAtom("!cat/pkg")
	assert.Nil(t, err)
	assert.Equal(t, atom.Blocker(), BlockerWeak)
	assert.Equal(t, atom.String(), "!cat/pkg")

	// subslotted
	atom, err = NewAtom("cat/pkg:1/2")
	assert.Nil(t, err)
	assert.Equal(t, atom.Slot(), "1")
	assert.Equal(t, atom.Subslot(), "2")
	assert.Equal(t, atom.String(), "cat/pkg:1/2")

	// slot operator
	atom, err = NewAtom("cat/pkg:0=")
	assert.Nil(t, err)
	assert.Equal(t, atom.Slot(), "0")
	assert.Equal(t, atom.SlotOp(), SlotOpEqual)
	assert.Equal(t, atom.String(), "cat/pkg:0=")

	// repo
	atom, err = NewAtom("cat/pkg::repo")
	assert.Nil(t, err)
	assert.Equal(t, atom.Repo(), "repo")
	assert.Equal(t, atom.String(), "cat/pkg::repo")

	// repo dep invalid on official EAPIs
	atom, err = NewAtomWithEapi("cat/pkg::repo", EAPI_LATEST)
	assert.Nil(t, atom)
	assert.NotNil(t, err)

	// all fields
	atom, err = NewAtom("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.Nil(t, err)
	assert.Equal(t, atom.Category(), "cat")
	assert.Equal(t, atom.Package(), "pkg")
	ver, _ = NewVersionWithOp("=1-r2")
	assert.Equal(t, atom.Version(), ver)
	assert.Equal(t, atom.Revision(), "2")
	assert.Equal(t, atom.Blocker(), BlockerStrong)
	assert.Equal(t, atom.Slot(), "3")
	assert.Equal(t, atom.Subslot(), "4")
	assert.Equal(t, atom.SlotOp(), SlotOpEqual)
	assert.Equal(t, atom.Use(), []string{"a", "b", "c"})
	assert.Equal(t, atom.Repo(), "repo")
	assert.Equal(t, atom.Cpn(), "cat/pkg")
	assert.Equal(t, atom.CPV(), "cat/pkg-1-r2")
	assert.Equal(t, atom.String(), "!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")

	// verify cached atoms reuse objects
	c1, _ = NewAtomCached("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.Equal(t, atom.Cmp(c1), 0)
	assert.True(t, atom != c1)
	c2, _ = NewAtomCached("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.True(t, c1 == c2)
	c1, _ = NewAtomCachedWithEapi("!!=a/b-1-r2:3/4=[a,b,c]", EAPI_LATEST)
	c2, _ = NewAtomCachedWithEapi("!!=a/b-1-r2:3/4=[a,b,c]", EAPI_LATEST)
	assert.True(t, c1 == c2)

	// a1 < a2
	a1, _ := NewAtom("=cat/pkg-1")
	a2, _ := NewAtom("=cat/pkg-2")
	assert.Equal(t, a1.Cmp(a2), -1)

	// a1 == a2
	a1, _ = NewAtom("=cat/pkg-2")
	a2, _ = NewAtom("=cat/pkg-2")
	assert.Equal(t, a1.Cmp(a2), 0)

	// a1 > a2
	a1, _ = NewAtom("=cat/pkg-2")
	a2, _ = NewAtom("=cat/pkg-1")
	assert.Equal(t, a1.Cmp(a2), 1)

	// hashing equal values
	a1, _ = NewAtom("=cat/pkg-1.0.2")
	a2, _ = NewAtom("=cat/pkg-1.0.2-r0")
	a3, _ := NewAtom("=cat/pkg-1.000.2")
	a4, _ := NewAtom("=cat/pkg-1.00.2-r0")
	m := make(map[uint64]bool)
	m[a1.Hash()] = true
	m[a2.Hash()] = true
	m[a3.Hash()] = true
	m[a4.Hash()] = true
	assert.Equal(t, len(m), 1)

	// hashing unequal values
	a1, _ = NewAtom("=cat/pkg-0.1")
	a2, _ = NewAtom("=cat/pkg-0.01")
	a3, _ = NewAtom("=cat/pkg-0.001")
	m = make(map[uint64]bool)
	m[a1.Hash()] = true
	m[a2.Hash()] = true
	m[a3.Hash()] = true
	assert.Equal(t, len(m), 3)
}

type validAtom struct {
	Atom     string
	Eapis    string
	Category string
	Package  string
	Blocker  string
	Version  string
	Revision string
	Slot     string
	Subslot  string
	Slot_Op  string
	Use      []string
}

type atomData struct {
	Valid   []validAtom
	Invalid [][]string
	Sorting [][][]string
}

func TestAtomToml(t *testing.T) {
	var atom_data atomData
	f, err := os.ReadFile("testdata/toml/atom.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(f, &atom_data)
	if err != nil {
		panic(err)
	}

	// valid atoms
	var ver *Version
	var blocker Blocker
	var slot_op SlotOperator
	for _, el := range atom_data.Valid {
		eapis, err := EapiRange(el.Eapis)
		if err != nil {
			panic(err)
		}
		for _, eapi := range eapis {
			atom, err := NewAtomWithEapi(el.Atom, eapi)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, atom.Category(), el.Category)
			assert.Equal(t, atom.Package(), el.Package)
			if el.Blocker == "" {
				blocker = BlockerNone
			} else {
				blocker, _ = BlockerFromString(el.Blocker)
			}
			assert.Equal(t, atom.Blocker(), blocker, "unequal blocker: %s", el.Atom)
			if el.Version != "" {
				ver, _ = NewVersionWithOp(el.Version)
			} else {
				ver = &Version{}
			}
			assert.Equal(t, atom.Version(), ver)
			assert.Equal(t, atom.Revision(), el.Revision)
			assert.Equal(t, atom.Slot(), el.Slot)
			assert.Equal(t, atom.Subslot(), el.Subslot)
			if el.Slot_Op == "" {
				slot_op = SlotOpNone
			} else {
				slot_op, _ = SlotOperatorFromString(el.Slot_Op)
			}
			assert.Equal(t, atom.SlotOp(), slot_op, "unequal slot ops: %s", el.Atom)
			assert.Equal(t, atom.Use(), el.Use, "unequal use: %s", el.Use)
		}
	}

	// invalid atoms
	for _, data := range atom_data.Invalid {
		s := data[0]
		failing_eapis, _ := EapiRange(data[1])
		failing_map := make(map[string]*Eapi)
		for _, eapi := range failing_eapis {
			failing_map[eapi.String()] = eapi
		}
		for _, eapi := range EAPIS {
			if _, ok := failing_map[eapi.String()]; ok {
				_, err := NewAtomWithEapi(s, eapi)
				assert.NotNil(t, err, "%s passed for EAPI=%s", s, eapi)
			} else {
				_, err := NewAtomWithEapi(s, eapi)
				assert.Nil(t, err, "%s failed for EAPI=%s", s, eapi)
			}
		}
	}

	// sorting
	for _, data := range atom_data.Sorting {
		var sorted []*Atom
		for _, s := range data[0] {
			atom, _ := NewAtom(s)
			sorted = append(sorted, atom)
		}
		sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Cmp(sorted[j]) == -1 })

		var expected []*Atom
		for _, s := range data[1] {
			atom, _ := NewAtom(s)
			expected = append(expected, atom)
		}

		assert.Equal(t, len(sorted), len(expected))
		for i := range sorted {
			assert.True(t, sorted[i].Cmp(expected[i]) == 0, "%s != %s", sorted, expected)
		}
	}
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
		assert.Equal(t, atom.String(), s)
	}
}

func BenchmarkNewAtom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		atom, _ := NewAtom(fmt.Sprintf("=cat/pkg-%d-r2:3/4=[a,b,c]", i))
		assert.NotNil(b, atom)
	}
}

func BenchmarkNewAtomStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		atom, _ := NewAtom("=cat/pkg-1-r2:3/4=[a,b,c]")
		assert.NotNil(b, atom)
	}
}

func BenchmarkNewAtomCachedStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		atom, _ := NewAtomCached("=cat/pkg-1-r2:3/4=[a,b,c]")
		assert.NotNil(b, atom)
	}
}

func BenchmarkAtomSort(b *testing.B) {
	var atoms []*Atom
	for i := 100; i > 0; i-- {
		a, _ := NewAtom(fmt.Sprintf("=cat/pkg-%d", i))
		atoms = append(atoms, a)
	}
	assert.Equal(b, atoms[0].String(), "=cat/pkg-100")
	for i := 0; i < b.N; i++ {
		sort.SliceStable(atoms, func(i, j int) bool { return atoms[i].Cmp(atoms[j]) == -1 })
	}
	assert.Equal(b, atoms[0].String(), "=cat/pkg-1")
}
