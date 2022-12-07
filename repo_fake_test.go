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

func TestExtend(t *testing.T) {
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
