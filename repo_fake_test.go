package pkgcraft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go"
)

func TestNewFakeRepo(t *testing.T) {
	// empty
	repo, _ := NewFakeRepo("fake", 0, []string{})
	assert.Equal(t, repo.Len(), 0)

	// single pkg
	repo, _ = NewFakeRepo("fake", 0, []string{"cat/pkg-1"})
	assert.Equal(t, repo.Len(), 1)

	// multiple pkgs with invalid cpv
	repo, _ = NewFakeRepo("fake", 0, []string{"a/b-1", "c/d-2", "=cat/pkg-1"})
	assert.Equal(t, repo.Len(), 2)
}

func TestFakeRepoExtend(t *testing.T) {
	// empty
	repo, _ := NewFakeRepo("fake", 0, []string{})
	assert.Equal(t, repo.Len(), 0)

	// add single pkg
	err := repo.Extend([]string{"cat/pkg-1"})
	assert.Equal(t, repo.Len(), 1)
	assert.Nil(t, err)

	// add multiple pkgs with overlap
	err = repo.Extend([]string{"cat/pkg-1", "cat/pkg-2"})
	assert.Equal(t, repo.Len(), 2)
	assert.Nil(t, err)
}

func TestFakeRepoPkgIter(t *testing.T) {
	var cpvs []string

	// empty
	repo, _ := NewFakeRepo("fake", 0, []string{})
	iter := repo.PkgIter()
	assert.False(t, iter.HasNext())
	assert.Nil(t, iter.Next())

	// add single pkg
	err := repo.Extend([]string{"cat/pkg-1"})
	assert.Nil(t, err)
	for iter := repo.PkgIter(); iter.HasNext(); {
		pkg := iter.Next()
		// verify repos are equal
		assert.True(t, repo.Cmp(pkg.Repo()) == 0)
		cpvs = append(cpvs, pkg.Atom().String())
	}
	assert.Equal(t, cpvs, []string{"cat/pkg-1"})

	// reset slice
	cpvs = cpvs[:0]

	// add multiple pkgs with overlap
	err = repo.Extend([]string{"cat/pkg-1", "cat/pkg-2"})
	assert.Nil(t, err)
	for iter := repo.PkgIter(); iter.HasNext(); {
		pkg := iter.Next()
		// verify repos are equal
		assert.True(t, repo.Cmp(pkg.Repo()) == 0)
		cpvs = append(cpvs, pkg.Atom().String())
	}
	assert.Equal(t, cpvs, []string{"cat/pkg-1", "cat/pkg-2"})
}

func TestFakeRepoPkgs(t *testing.T) {
	var cpvs []string

	// empty
	repo, _ := NewFakeRepo("fake", 0, []string{})
	assert.Empty(t, repo.Pkgs())

	// add single pkg
	err := repo.Extend([]string{"cat/pkg-1"})
	assert.Nil(t, err)
	for pkg := range repo.Pkgs() {
		// verify repos are equal
		assert.True(t, repo.Cmp(pkg.Repo()) == 0)
		cpvs = append(cpvs, pkg.Atom().String())
	}
	assert.Equal(t, cpvs, []string{"cat/pkg-1"})

	// reset slice
	cpvs = cpvs[:0]

	// add multiple pkgs with overlap
	err = repo.Extend([]string{"cat/pkg-1", "cat/pkg-2"})
	assert.Nil(t, err)
	for pkg := range repo.Pkgs() {
		// verify repos are equal
		assert.True(t, repo.Cmp(pkg.Repo()) == 0)
		cpvs = append(cpvs, pkg.Atom().String())
	}
	assert.Equal(t, cpvs, []string{"cat/pkg-1", "cat/pkg-2"})
}

func TestFakeRepoRestrictPkgIter(t *testing.T) {
	var cpvs []string
	restrict, _ := NewRestrict("<cat/pkg-2")

	// empty
	repo, _ := NewFakeRepo("fake", 0, []string{})
	iter := repo.RestrictPkgIter(restrict)
	assert.False(t, iter.HasNext())
	assert.Nil(t, iter.Next())

	// add single pkg
	err := repo.Extend([]string{"cat/pkg-1"})
	assert.Nil(t, err)
	for iter := repo.RestrictPkgIter(restrict); iter.HasNext(); {
		pkg := iter.Next()
		// verify repos are equal
		assert.True(t, repo.Cmp(pkg.Repo()) == 0)
		cpvs = append(cpvs, pkg.Atom().String())
	}
	assert.Equal(t, cpvs, []string{"cat/pkg-1"})

	// reset slice
	cpvs = cpvs[:0]

	// add multiple pkgs with overlap
	err = repo.Extend([]string{"cat/pkg-1", "cat/pkg-2"})
	assert.Nil(t, err)
	for iter := repo.RestrictPkgIter(restrict); iter.HasNext(); {
		pkg := iter.Next()
		// verify repos are equal
		assert.True(t, repo.Cmp(pkg.Repo()) == 0)
		cpvs = append(cpvs, pkg.Atom().String())
	}
	assert.Equal(t, cpvs, []string{"cat/pkg-1"})
}

func TestFakeRepoRestrictPkgs(t *testing.T) {
	var cpvs []string
	restrict, _ := NewRestrict("<cat/pkg-2")

	// empty
	repo, _ := NewFakeRepo("fake", 0, []string{})
	assert.Empty(t, repo.RestrictPkgs(restrict))

	// add single pkg
	err := repo.Extend([]string{"cat/pkg-1"})
	assert.Nil(t, err)
	for pkg := range repo.RestrictPkgs(restrict) {
		// verify repos are equal
		assert.True(t, repo.Cmp(pkg.Repo()) == 0)
		cpvs = append(cpvs, pkg.Atom().String())
	}
	assert.Equal(t, cpvs, []string{"cat/pkg-1"})

	// reset slice
	cpvs = cpvs[:0]

	// add multiple pkgs with overlap
	err = repo.Extend([]string{"cat/pkg-1", "cat/pkg-2"})
	assert.Nil(t, err)
	for pkg := range repo.RestrictPkgs(restrict) {
		// verify repos are equal
		assert.True(t, repo.Cmp(pkg.Repo()) == 0)
		cpvs = append(cpvs, pkg.Atom().String())
	}
	assert.Equal(t, cpvs, []string{"cat/pkg-1"})
}
