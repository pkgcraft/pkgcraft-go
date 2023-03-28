package pkgcraft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pkgcraft/pkgcraft-go/internal"
)

func TestReverseSlice(t *testing.T) {
	// empty
	nums := []int{}
	ReverseSlice(nums)
	assert.Equal(t, nums, []int{})

	// single value
	nums = []int{1}
	ReverseSlice(nums)
	assert.Equal(t, nums, []int{1})

	// multiple values
	nums = []int{1, 2, 3}
	ReverseSlice(nums)
	assert.Equal(t, nums, []int{3, 2, 1})
}
