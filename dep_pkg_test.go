package pkgcraft_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
	. "github.com/pkgcraft/pkgcraft-go/internal"
)

func TestBlockerFromString(t *testing.T) {
	valid := map[string]Blocker{"!": BlockerWeak, "!!": BlockerStrong}
	for s, expected := range valid {
		v, err := BlockerFromString(s)
		assert.Equal(t, v, expected)
		assert.Nil(t, err)
	}

	invalid := []string{"", "!!!", "a"}
	for _, s := range invalid {
		_, err := BlockerFromString(s)
		assert.NotNil(t, err)
	}
}

func TestSlotOperatorFromString(t *testing.T) {
	valid := map[string]SlotOperator{"=": SlotOpEqual, "*": SlotOpStar}
	for s, expected := range valid {
		v, err := SlotOperatorFromString(s)
		assert.Equal(t, v, expected)
		assert.Nil(t, err)
	}

	invalid := []string{"", "=*", "*=", "~"}
	for _, s := range invalid {
		_, err := SlotOperatorFromString(s)
		assert.NotNil(t, err)
	}
}

func TestDepAttrs(t *testing.T) {
	var dep, c1, c2 *Dep
	var err error
	var ver *Version

	// unversioned
	dep, err = NewDep("cat/pkg")
	assert.Nil(t, err)
	assert.Equal(t, dep.Category(), "cat")
	assert.Equal(t, dep.Package(), "pkg")
	assert.Equal(t, dep.Version(), &Version{})
	assert.Equal(t, dep.Revision(), "")
	assert.Equal(t, dep.P(), "pkg")
	assert.Equal(t, dep.Pf(), "pkg")
	assert.Equal(t, dep.Pr(), "")
	assert.Equal(t, dep.Pv(), "")
	assert.Equal(t, dep.Pvr(), "")
	assert.Equal(t, dep.Cpn(), "cat/pkg")
	assert.Equal(t, dep.Cpv(), "cat/pkg")
	assert.Equal(t, dep.Blocker(), BlockerNone)
	assert.Equal(t, dep.Slot(), "")
	assert.Equal(t, dep.Subslot(), "")
	assert.Equal(t, dep.SlotOp(), SlotOpNone)
	assert.Equal(t, len(dep.Use()), 0)
	assert.Equal(t, dep.Repo(), "")
	assert.Equal(t, dep.String(), "cat/pkg")

	// versioned
	dep, err = NewDep("=cat/pkg-1-r2")
	assert.Nil(t, err)
	assert.Equal(t, dep.Category(), "cat")
	assert.Equal(t, dep.Package(), "pkg")
	ver, _ = NewVersion("=1-r2")
	assert.True(t, dep.Version().Cmp(ver) == 0)
	assert.Equal(t, dep.Revision(), "2")
	assert.Equal(t, dep.P(), "pkg-1")
	assert.Equal(t, dep.Pf(), "pkg-1-r2")
	assert.Equal(t, dep.Pr(), "r2")
	assert.Equal(t, dep.Pv(), "1")
	assert.Equal(t, dep.Pvr(), "1-r2")
	assert.Equal(t, dep.Cpn(), "cat/pkg")
	assert.Equal(t, dep.Cpv(), "cat/pkg-1-r2")
	assert.Equal(t, dep.String(), "=cat/pkg-1-r2")

	// blocker
	dep, err = NewDep("!cat/pkg")
	assert.Nil(t, err)
	assert.Equal(t, dep.Blocker(), BlockerWeak)
	assert.Equal(t, dep.String(), "!cat/pkg")

	// subslotted
	dep, err = NewDep("cat/pkg:1/2")
	assert.Nil(t, err)
	assert.Equal(t, dep.Slot(), "1")
	assert.Equal(t, dep.Subslot(), "2")
	assert.Equal(t, dep.String(), "cat/pkg:1/2")

	// slot operator
	dep, err = NewDep("cat/pkg:0=")
	assert.Nil(t, err)
	assert.Equal(t, dep.Slot(), "0")
	assert.Equal(t, dep.SlotOp(), SlotOpEqual)
	assert.Equal(t, dep.String(), "cat/pkg:0=")

	// repo
	dep, err = NewDep("cat/pkg::repo")
	assert.Nil(t, err)
	assert.Equal(t, dep.Repo(), "repo")
	assert.Equal(t, dep.String(), "cat/pkg::repo")

	// repo dep invalid on official EAPIs
	dep, err = NewDepWithEapi("cat/pkg::repo", EAPI_LATEST_OFFICIAL)
	assert.Nil(t, dep)
	assert.NotNil(t, err)

	// all fields
	dep, err = NewDep("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.Nil(t, err)
	assert.Equal(t, dep.Category(), "cat")
	assert.Equal(t, dep.Package(), "pkg")
	ver, _ = NewVersion("=1-r2")
	assert.True(t, dep.Version().Cmp(ver) == 0)
	assert.Equal(t, dep.Revision(), "2")
	assert.Equal(t, dep.Blocker(), BlockerStrong)
	assert.Equal(t, dep.Slot(), "3")
	assert.Equal(t, dep.Subslot(), "4")
	assert.Equal(t, dep.SlotOp(), SlotOpEqual)
	assert.Equal(t, dep.Use(), []string{"a", "b", "c"})
	assert.Equal(t, dep.Repo(), "repo")
	assert.Equal(t, dep.P(), "pkg-1")
	assert.Equal(t, dep.Pf(), "pkg-1-r2")
	assert.Equal(t, dep.Pr(), "r2")
	assert.Equal(t, dep.Pv(), "1")
	assert.Equal(t, dep.Pvr(), "1-r2")
	assert.Equal(t, dep.Cpn(), "cat/pkg")
	assert.Equal(t, dep.Cpv(), "cat/pkg-1-r2")
	assert.Equal(t, dep.String(), "!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")

	// verify cached deps reuse objects
	c1, _ = NewDepCached("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.Equal(t, dep.Cmp(c1), 0)
	assert.True(t, dep != c1)
	c2, _ = NewDepCached("!!=cat/pkg-1-r2:3/4=[a,b,c]::repo")
	assert.True(t, c1 == c2)
	c1, _ = NewDepCachedWithEapi("!!=a/b-1-r2:3/4=[a,b,c]", EAPI_LATEST_OFFICIAL)
	c2, _ = NewDepCachedWithEapi("!!=a/b-1-r2:3/4=[a,b,c]", EAPI_LATEST_OFFICIAL)
	assert.True(t, c1 == c2)
}

func TestDepCmp(t *testing.T) {
	// d1 < d2
	d1, _ := NewDep("=cat/pkg-1")
	d2, _ := NewDep("=cat/pkg-2")
	assert.Equal(t, d1.Cmp(d2), -1)

	// d1 == d2
	d1, _ = NewDep("=cat/pkg-2")
	d2, _ = NewDep("=cat/pkg-2")
	assert.Equal(t, d1.Cmp(d2), 0)

	// d1 > d2
	d1, _ = NewDep("=cat/pkg-2")
	d2, _ = NewDep("=cat/pkg-1")
	assert.Equal(t, d1.Cmp(d2), 1)
}

func TestDepHash(t *testing.T) {
	for _, data := range VERSION_TOML.Hashing {
		m := make(map[uint64]bool)
		for _, s := range data.Versions {
			dep, _ := NewDep(fmt.Sprintf("=cat/pkg-%s", s))
			m[dep.Hash()] = true
		}

		if data.Equal {
			assert.Equal(t, len(m), 1)
		} else {
			assert.Equal(t, len(m), len(data.Versions))
		}
	}
}

// TODO: use shared intersects test data
func TestDepIntersects(t *testing.T) {
	// equal versions
	d1, _ := NewDep("=a/b-1.0.2")
	d2, _ := NewDep("=a/b-1.0.2-r0")
	assert.True(t, d1.Intersects(d2))

	// unequal versions
	d1, _ = NewDep("=a/b-0")
	d2, _ = NewDep("=a/b-1")
	assert.False(t, d1.Intersects(d2))

	// CPV and dep
	cpv, _ := NewCpv("a/b-0")
	dep, _ := NewDep("=a/b-0*")
	assert.True(t, cpv.Intersects(dep))
}

func TestDepParse(t *testing.T) {
	// valid
	var ver *Version
	var blocker Blocker
	var slot_op SlotOperator
	for _, el := range DEP_TOML.Valid {
		eapis, err := EapiRange(el.Eapis)
		if err != nil {
			panic(err)
		}
		for _, eapi := range eapis {
			dep, err := NewDepWithEapi(el.Dep, eapi)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, dep.Category(), el.Category)
			assert.Equal(t, dep.Package(), el.Package)
			if el.Blocker == "" {
				blocker = BlockerNone
			} else {
				blocker, _ = BlockerFromString(el.Blocker)
			}
			assert.Equal(t, dep.Blocker(), blocker, "unequal blocker: %s", el.Dep)
			if el.Version != "" {
				ver, _ = NewVersion(el.Version)
				assert.True(t, dep.Version().Cmp(ver) == 0)
			} else {
				assert.Equal(t, dep.Version(), &Version{})
			}
			assert.Equal(t, dep.Revision(), el.Revision)
			assert.Equal(t, dep.Slot(), el.Slot)
			assert.Equal(t, dep.Subslot(), el.Subslot)
			if el.Slot_Op == "" {
				slot_op = SlotOpNone
			} else {
				slot_op, _ = SlotOperatorFromString(el.Slot_Op)
			}
			assert.Equal(t, dep.SlotOp(), slot_op, "unequal slot ops: %s", el.Dep)
			assert.Equal(t, dep.Use(), el.Use, "unequal use: %s", el.Use)
		}
	}

	// invalid
	for _, s := range DEP_TOML.Invalid {
		for _, eapi := range EAPIS {
			_, err := NewDepWithEapi(s, eapi)
			assert.NotNil(t, err, "%s passed for EAPI=%s", s, eapi)
		}
	}
}

func TestDepSort(t *testing.T) {
	for _, data := range DEP_TOML.Sorting {
		var expected []*Dep
		for _, s := range data.Sorted {
			dep, _ := NewDep(s)
			expected = append(expected, dep)
		}

		sorted := make([]*Dep, len(expected))
		copy(sorted, expected)
		ReverseSlice(sorted)
		sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Cmp(sorted[j]) == -1 })

		// equal deps aren't sorted so reversing should restore the original order
		if data.Equal {
			ReverseSlice(sorted)
		}

		assert.Equal(t, len(sorted), len(expected))
		for i := range sorted {
			assert.True(t, sorted[i].Cmp(expected[i]) == 0, "%s != %s", sorted, expected)
		}
	}
}

// test sending Deps over a channel
func TestDepChannels(t *testing.T) {
	var dep *Dep

	dep_strs := make(chan string)
	deps := make(chan *Dep)

	go func() {
		for {
			s := <-dep_strs
			dep, _ = NewDep(s)
			deps <- dep
		}
	}()

	var s string
	for i := 0; i < 1000; i++ {
		s = fmt.Sprintf("=cat/pkg-%d-r2:3/4=[a,b,c]", i)
		dep_strs <- s
		dep = <-deps
		assert.Equal(t, dep.String(), s)
	}
}

func BenchmarkNewDep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dep, _ := NewDep(fmt.Sprintf("=cat/pkg-%d-r2:3/4=[a,b,c]", i))
		assert.NotNil(b, dep)
	}
}

func BenchmarkNewDepStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dep, _ := NewDep("=cat/pkg-1-r2:3/4=[a,b,c]")
		assert.NotNil(b, dep)
	}
}

func BenchmarkNewDepCachedStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dep, _ := NewDepCached("=cat/pkg-1-r2:3/4=[a,b,c]")
		assert.NotNil(b, dep)
	}
}

func BenchmarkDepSort(b *testing.B) {
	var deps []*Dep
	for i := 100; i > 0; i-- {
		a, _ := NewDep(fmt.Sprintf("=cat/pkg-%d", i))
		deps = append(deps, a)
	}
	assert.Equal(b, deps[0].String(), "=cat/pkg-100")
	for i := 0; i < b.N; i++ {
		sort.SliceStable(deps, func(i, j int) bool { return deps[i].Cmp(deps[j]) == -1 })
	}
	assert.Equal(b, deps[0].String(), "=cat/pkg-1")
}
