package pkgcraft_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestVersion(t *testing.T) {
	var version *Version

	// non-revision
	version, _ = NewVersion("1")
	assert.Equal(t, version.Revision(), "0")
	assert.Equal(t, fmt.Sprintf("%s", version), "1")

	// revisioned
	version, _ = NewVersion("1-r1")
	assert.Equal(t, version.Revision(), "1")
	assert.Equal(t, fmt.Sprintf("%s", version), "1-r1")

	// explicit '0' revision
	version, _ = NewVersion("1-r0")
	assert.Equal(t, version.Revision(), "0")
	assert.Equal(t, fmt.Sprintf("%s", version), "1-r0")

	// invalid
	version, _ = NewVersion(">1-r2")
	assert.Nil(t, version)

	// Version with op
	version, _ = NewVersionWithOp(">1-r2")
	assert.Equal(t, version.Revision(), "2")
	assert.Equal(t, fmt.Sprintf("%s", version), ">1-r2")

	// v1 < v2
	v1, _ := NewVersion("1")
	v2, _ := NewVersion("2")
	assert.Equal(t, v1.Cmp(v2), -1)

	// v1 == v2
	v1, _ = NewVersion("2")
	v2, _ = NewVersion("2")
	assert.Equal(t, v1.Cmp(v2), 0)

	// v1 > v2
	v1, _ = NewVersion("2")
	v2, _ = NewVersion("1")
	assert.Equal(t, v1.Cmp(v2), 1)

	// hashing equal values
	v1, _ = NewVersion("1.0.2")
	v2, _ = NewVersion("1.0.2-r0")
	v3, _ := NewVersion("1.000.2")
	v4, _ := NewVersion("1.00.2-r0")
	m := make(map[uint64]bool)
	m[v1.Hash()] = true
	m[v2.Hash()] = true
	m[v3.Hash()] = true
	m[v4.Hash()] = true
	assert.Equal(t, len(m), 1)

	// hashing unequal values
	v1, _ = NewVersion("0.1")
	v2, _ = NewVersion("0.01")
	v3, _ = NewVersion("0.001")
	m = make(map[uint64]bool)
	m[v1.Hash()] = true
	m[v2.Hash()] = true
	m[v3.Hash()] = true
	assert.Equal(t, len(m), 3)
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
