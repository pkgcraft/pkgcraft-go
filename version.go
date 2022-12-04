package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type Version struct {
	ptr *C.AtomVersion
}

func newVersion(ptr *C.AtomVersion) (*Version, error) {
	if ptr != nil {
		ver := &Version{ptr}
		runtime.SetFinalizer(ver, func(v *Version) { C.pkgcraft_version_free(v.ptr) })
		return ver, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

// Parse a string into a version.
func NewVersion(s string) (*Version, error) {
	ver_str := C.CString(s)
	defer C.free(unsafe.Pointer(ver_str))
	ptr := C.pkgcraft_version_new(ver_str)
	return newVersion(ptr)
}

// Parse a string into a version with an operator.
func NewVersionWithOp(s string) (*Version, error) {
	ver_str := C.CString(s)
	defer C.free(unsafe.Pointer(ver_str))
	ptr := C.pkgcraft_version_with_op(ver_str)
	return newVersion(ptr)
}

// Return a version's revision.
func (v *Version) Revision() string {
	s := C.pkgcraft_version_revision(v.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Compare a version with another version returning -1, 0, or 1 if the first
// version is less than, equal to, or greater than the second version,
// respectively.
func (v1 *Version) Cmp(v2 *Version) int {
	return int(C.pkgcraft_version_cmp(v1.ptr, v2.ptr))
}

func (v *Version) String() string {
	s := C.pkgcraft_version_str(v.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (v *Version) Hash() uint64 {
	return uint64(C.pkgcraft_version_hash(v.ptr))
}
