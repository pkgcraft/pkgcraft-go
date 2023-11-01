package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
	"unsafe"
)

type Revision struct {
	ptr *C.Revision
}

type revisionPtr interface {
	p() *C.Revision
}

func revisionFromPtr(ptr *C.Revision) (*Revision, error) {
	if ptr != nil {
		rev := &Revision{ptr}
		runtime.SetFinalizer(rev, func(self *Revision) { C.pkgcraft_revision_free(self.ptr) })
		return rev, nil
	} else {
		return nil, newPkgcraftError()
	}
}

// Parse a string into a revision.
func NewRevision(s string) (*Revision, error) {
	ver_str := C.CString(s)
	defer C.free(unsafe.Pointer(ver_str))
	ptr := C.pkgcraft_revision_new(ver_str)
	return revisionFromPtr(ptr)
}

func (self *Revision) p() *C.Revision {
	return self.ptr
}

// Compare a revision with another revision returning -1, 0, or 1 if the first
// is less than, equal to, or greater than the second, respectively.
func (self *Revision) Cmp(other revisionPtr) int {
	return int(C.pkgcraft_revision_cmp(self.ptr, other.p()))
}

func (self *Revision) String() string {
	s := C.pkgcraft_revision_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Revision) Hash() uint64 {
	return uint64(C.pkgcraft_revision_hash(self.ptr))
}

type Version struct {
	ptr *C.Version
	// cached fields
	_revision *Revision
}

type versionPtr interface {
	p() *C.Version
}

func versionFromPtr(ptr *C.Version) (*Version, error) {
	if ptr != nil {
		ver := &Version{ptr: ptr}
		runtime.SetFinalizer(ver, func(self *Version) { C.pkgcraft_version_free(self.ptr) })
		return ver, nil
	} else {
		return nil, newPkgcraftError()
	}
}

// Parse a string into a version.
func NewVersion(s string) (*Version, error) {
	ver_str := C.CString(s)
	defer C.free(unsafe.Pointer(ver_str))
	ptr := C.pkgcraft_version_new(ver_str)
	return versionFromPtr(ptr)
}

func (self *Version) p() *C.Version {
	return self.ptr
}

// Return a version's revision.
func (self *Version) Revision() *Revision {
	if self._revision == nil {
		ptr := C.pkgcraft_version_revision(self.ptr)
		if ptr != nil {
			self._revision, _ = revisionFromPtr(ptr)
		} else {
			self._revision = &Revision{}
		}
	}
	return self._revision
}

// Compare a version with another version returning -1, 0, or 1 if the first is
// less than, equal to, or greater than the second, respectively.
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
