package pkgcraft

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLibVersion(t *testing.T) {
	ver, _ := version.NewVersion(pkgcraft_lib_version())
	min_ver, _ := version.NewVersion("0.0.2")
	max_ver, _ := version.NewVersion("0.0.2")
	min_err := fmt.Sprintf("pkgcraft C library %s failed requirements >=%s\n", ver, min_ver)
	max_err := fmt.Sprintf("pkgcraft C library %s failed requirements <=%s\n", ver, max_ver)

	assert.True(t, ver.GreaterThanOrEqual(min_ver), min_err)
	assert.True(t, ver.LessThanOrEqual(max_ver), max_err)
}
