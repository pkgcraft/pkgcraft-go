package pkgcraft_test

import (
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
