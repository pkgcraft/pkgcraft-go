package pkgcraft

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

// version requirements for pkgcraft C library
const MIN_VERSION = "0.0.3"
const MAX_VERSION = "0.0.3"

func TestLibVersion(t *testing.T) {
	ver, _ := version.NewVersion(pkgcraftLibVersion())
	min_ver, _ := version.NewVersion(MIN_VERSION)
	max_ver, _ := version.NewVersion(MAX_VERSION)
	min_err := fmt.Sprintf("pkgcraft C library %s failed requirements >=%s\n", ver, min_ver)
	max_err := fmt.Sprintf("pkgcraft C library %s failed requirements <=%s\n", ver, max_ver)

	assert.True(t, ver.GreaterThanOrEqual(min_ver), min_err)
	assert.True(t, ver.LessThanOrEqual(max_ver), max_err)
}
