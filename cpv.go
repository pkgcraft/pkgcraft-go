package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
	"unsafe"
)

type Cpv struct {
	ptr *C.Cpv
	// cached fields
	_category string
	_package  string
	_version  *Version
	_hash     uint64
}

type cpvPtr interface {
	p() *C.Cpv
}

func cpvFromPtr(ptr *C.Cpv) (*Cpv, error) {
	if ptr != nil {
		cpv := &Cpv{ptr: ptr}
		runtime.SetFinalizer(cpv, func(self *Cpv) { C.pkgcraft_cpv_free(self.ptr) })
		return cpv, nil
	} else {
		return nil, newPkgcraftError()
	}
}

// Parse a string into a Cpv object.
func NewCpv(s string) (*Cpv, error) {
	c_str := C.CString(s)
	defer C.free(unsafe.Pointer(c_str))
	ptr := C.pkgcraft_cpv_new(c_str)
	return cpvFromPtr(ptr)
}

func (self *Cpv) p() *C.Cpv {
	return self.ptr
}

// Return an Cpv's category.
func (self *Cpv) Category() string {
	if self._category == "" {
		s := C.pkgcraft_cpv_category(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._category = C.GoString(s)
	}
	return self._category
}

// Return a Cpv's package name.
func (self *Cpv) Package() string {
	if self._package == "" {
		s := C.pkgcraft_cpv_package(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._package = C.GoString(s)
	}
	return self._package
}

// Return a Cpv's version.
func (self *Cpv) Version() *Version {
	if self._version == nil {
		ptr := C.pkgcraft_cpv_version(self.ptr)
		if ptr != nil {
			self._version, _ = versionFromPtr(ptr)
		} else {
			self._version = &Version{}
		}
	}
	return self._version
}

// Return a Cpv's revision.
func (self *Cpv) Revision() string {
	version := self.Version()
	if *version != (Version{}) {
		return version.Revision()
	}
	return ""
}

// Return a Cpv's package and version.
func (self *Cpv) P() string {
	c_str := C.pkgcraft_cpv_p(self.ptr)
	defer C.pkgcraft_str_free(c_str)
	return C.GoString(c_str)
}

// Return a Cpv's package, version, and revision.
func (self *Cpv) Pf() string {
	c_str := C.pkgcraft_cpv_pf(self.ptr)
	defer C.pkgcraft_str_free(c_str)
	return C.GoString(c_str)
}

// Return a Cpv's revision.
func (self *Cpv) Pr() string {
	c_str := C.pkgcraft_cpv_pr(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return a Cpv's version.
func (self *Cpv) Pv() string {
	c_str := C.pkgcraft_cpv_pv(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return a Cpv's version and revision.
func (self *Cpv) Pvr() string {
	c_str := C.pkgcraft_cpv_pvr(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return a Cpv's category and package.
func (self *Cpv) Cpn() string {
	s := C.pkgcraft_cpv_cpn(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Cpv) String() string {
	s := C.pkgcraft_cpv_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Cpv) Hash() uint64 {
	if self._hash == 0 {
		self._hash = uint64(C.pkgcraft_cpv_hash(self.ptr))
	}
	return self._hash
}

// Compare two deps returning -1, 0, or 1 if the first is less than, equal to,
// or greater than the second, respectively.
func (self *Cpv) Cmp(other *Cpv) int {
	return int(C.pkgcraft_cpv_cmp(self.ptr, other.ptr))
}

// Determine if two Cpv or Dep objects intersect.
func (self *Cpv) Intersects(other interface{}) bool {
	switch other := other.(type) {
	case *Cpv:
		return bool(C.pkgcraft_cpv_intersects(self.ptr, other.ptr))
	case *Dep:
		return bool(C.pkgcraft_cpv_intersects_dep(self.ptr, other.ptr))
	default:
		return false
	}
}
