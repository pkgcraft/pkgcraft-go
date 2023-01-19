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

type versionPtr interface {
	p() *C.AtomVersion
}

func versionFromPtr(ptr *C.AtomVersion) (*Version, error) {
	if ptr != nil {
		ver := &Version{ptr}
		runtime.SetFinalizer(ver, func(self *Version) { C.pkgcraft_version_free(self.ptr) })
		return ver, nil
	} else {
		err := C.pkgcraft_error_last()
		defer C.pkgcraft_error_free(err)
		return nil, errors.New(C.GoString(err.message))
	}
}

// Parse a string into a version.
func NewVersion(s string) (*Version, error) {
	ver_str := C.CString(s)
	defer C.free(unsafe.Pointer(ver_str))
	ptr := C.pkgcraft_version_new(ver_str)
	return versionFromPtr(ptr)
}

func (self *Version) p() *C.AtomVersion {
	return self.ptr
}

// Return a version's revision.
func (self *Version) Revision() string {
	s := C.pkgcraft_version_revision(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Compare a version with another version returning -1, 0, or 1 if the first
// version is less than, equal to, or greater than the second version,
// respectively.
func (self *Version) Cmp(other versionPtr) int {
	return int(C.pkgcraft_version_cmp(self.ptr, other.p()))
}

func (self *Version) String() string {
	s := C.pkgcraft_version_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Version) Hash() uint64 {
	return uint64(C.pkgcraft_version_hash(self.ptr))
}

// Determine if two versions intersect.
func (self *Version) Intersects(other versionPtr) bool {
	return bool(C.pkgcraft_version_intersects(self.ptr, other.p()))
}

type VersionWithOp struct {
	*Version
}

// Parse a string into a version with an operator.
func NewVersionWithOp(s string) (*VersionWithOp, error) {
	ver_str := C.CString(s)
	defer C.free(unsafe.Pointer(ver_str))
	ptr := C.pkgcraft_version_with_op(ver_str)
	ver, err := versionFromPtr(ptr)
	if ver != nil {
		return &VersionWithOp{ver}, nil
	} else {
		return nil, err
	}
}

func (self *VersionWithOp) String() string {
	s := C.pkgcraft_version_str_with_op(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
