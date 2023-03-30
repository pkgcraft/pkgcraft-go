package pkgcraft_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
	. "github.com/pkgcraft/pkgcraft-go/internal"
)

func TestVersion(t *testing.T) {
	// non-revision
	ver, err := NewVersion("1")
	assert.Nil(t, err)
	assert.Equal(t, ver.Revision(), "")
	assert.Equal(t, ver.String(), "1")

	// revisioned
	ver, err = NewVersion("1-r1")
	assert.Nil(t, err)
	assert.Equal(t, ver.Revision(), "1")
	assert.Equal(t, ver.String(), "1-r1")

	// explicit '0' revision
	ver, err = NewVersion("1-r0")
	assert.Nil(t, err)
	assert.Equal(t, ver.Revision(), "0")
	assert.Equal(t, ver.String(), "1-r0")

	// invalid
	ver, err = NewVersion(">1-r2")
	assert.Nil(t, ver)
	assert.NotNil(t, err)
}

func TestVersionWithOp(t *testing.T) {
	ver, err := NewVersionWithOp(">1-r2")
	assert.Nil(t, err)
	assert.Equal(t, ver.Revision(), "2")
	assert.Equal(t, ver.String(), ">1-r2")
}

func TestVersionCmp(t *testing.T) {
	for _, s := range VERSION_TOML.Compares {
		vals := strings.Fields(s)
		v1, _ := NewVersion(vals[0])
		op := vals[1]
		v2, _ := NewVersion(vals[2])

		switch op {
		case "<":
			assert.Equal(t, v1.Cmp(v2), -1)
		case "==":
			assert.Equal(t, v1.Cmp(v2), 0)
		case "!=":
			assert.NotEqual(t, v1.Cmp(v2), 0)
		case ">":
			assert.Equal(t, v1.Cmp(v2), 1)
		default:
			panic(fmt.Sprintf("invalid operator: %s", op))
		}
	}
}

func TestVersionHash(t *testing.T) {
	for _, data := range VERSION_TOML.Hashing {
		m := make(map[uint64]bool)
		for _, s := range data.Versions {
			ver, _ := NewVersion(s)
			m[ver.Hash()] = true
		}

		if data.Equal {
			assert.Equal(t, len(m), 1)
		} else {
			assert.Equal(t, len(m), len(data.Versions))
		}
	}
}

// TODO: use shared intersects test data
func TestVersionIntersects(t *testing.T) {
	var v1, v2 *Version
	var vo1, vo2 *VersionWithOp

	// equal, non-op versions
	v1, _ = NewVersion("1.0.2")
	v2, _ = NewVersion("1.0.2-r0")
	assert.True(t, v1.Intersects(v2))

	// unequal, non-op versions
	v1, _ = NewVersion("0")
	v2, _ = NewVersion("1")
	assert.False(t, v1.Intersects(v2))

	// non-op and op versions
	vo1, _ = NewVersionWithOp("<0")
	v2, _ = NewVersion("0")
	assert.False(t, vo1.Intersects(v2))
	v1, _ = NewVersion("0")
	vo2, _ = NewVersionWithOp("=0*")
	assert.True(t, v1.Intersects(vo2))
}

func TestVersionSort(t *testing.T) {
	for _, data := range VERSION_TOML.Sorting {
		var expected []*Version
		for _, s := range data.Sorted {
			ver, _ := NewVersion(s)
			expected = append(expected, ver)
		}

		sorted := make([]*Version, len(expected))
		copy(sorted, expected)
		ReverseSlice(sorted)
		sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Cmp(sorted[j]) == -1 })

		// equal versions aren't sorted so reversing should restore the original order
		if data.Equal {
			ReverseSlice(sorted)
		}

		assert.Equal(t, len(sorted), len(expected))
		for i := range sorted {
			assert.True(t, sorted[i].Cmp(expected[i]) == 0, "%s != %s", sorted, expected)
		}
	}
}

func BenchmarkNewVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ver, _ := NewVersion("1.2.3_alpha4-r5")
		assert.NotNil(b, ver)
	}
}

func BenchmarkVersionSort(b *testing.B) {
	var versions []*Version
	for i := 100; i > 0; i-- {
		v, _ := NewVersion(fmt.Sprintf("%d", i))
		versions = append(versions, v)
	}
	assert.Equal(b, versions[0].String(), "100")
	for i := 0; i < b.N; i++ {
		sort.SliceStable(versions, func(i, j int) bool { return versions[i].Cmp(versions[j]) == -1 })
	}
	assert.Equal(b, versions[0].String(), "1")
}
