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
