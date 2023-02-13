package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type Cpv struct {
	ptr *C.Dep
	// cached fields
	_category string
	_package  string
	_version  *Version
	_hash     uint64
}

type depPtr interface {
	p() *C.Dep
}

func cpvFromPtr(ptr *C.Dep) (*Cpv, error) {
	if ptr != nil {
		cpv := &Cpv{ptr: ptr}
		runtime.SetFinalizer(cpv, func(self *Cpv) { C.pkgcraft_dep_free(self.ptr) })
		return cpv, nil
	} else {
		err := C.pkgcraft_error_last()
		defer C.pkgcraft_error_free(err)
		return nil, errors.New(C.GoString(err.message))
	}
}

// Parse a CPV string into an dep.
func NewCpv(s string) (*Cpv, error) {
	c_str := C.CString(s)
	defer C.free(unsafe.Pointer(c_str))
	ptr := C.pkgcraft_cpv_new(c_str)
	return cpvFromPtr(ptr)
}

func (self *Cpv) p() *C.Dep {
	return self.ptr
}

// Return an dep's category.
func (self *Cpv) Category() string {
	if self._category == "" {
		s := C.pkgcraft_dep_category(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._category = C.GoString(s)
	}
	return self._category
}

// Return an dep's package name.
func (self *Cpv) Package() string {
	if self._package == "" {
		s := C.pkgcraft_dep_package(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._package = C.GoString(s)
	}
	return self._package
}

// Return an dep's version.
func (self *Cpv) Version() *Version {
	if self._version == nil {
		ptr := C.pkgcraft_dep_version(self.ptr)
		if ptr != nil {
			self._version, _ = versionFromPtr(ptr)
		} else {
			self._version = &Version{}
		}
	}
	return self._version
}

// Return an dep's revision.
func (self *Cpv) Revision() string {
	version := self.Version()
	if *version != (Version{}) {
		return version.Revision()
	}
	return ""
}

// Return an dep's package and version.
func (self *Cpv) P() string {
	c_str := C.pkgcraft_dep_p(self.ptr)
	defer C.pkgcraft_str_free(c_str)
	return C.GoString(c_str)
}

// Return an dep's package, version, and revision.
func (self *Cpv) Pf() string {
	c_str := C.pkgcraft_dep_pf(self.ptr)
	defer C.pkgcraft_str_free(c_str)
	return C.GoString(c_str)
}

// Return an dep's revision.
func (self *Cpv) Pr() string {
	c_str := C.pkgcraft_dep_pr(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return an dep's version.
func (self *Cpv) Pv() string {
	c_str := C.pkgcraft_dep_pv(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return an dep's version and revision.
func (self *Cpv) Pvr() string {
	c_str := C.pkgcraft_dep_pvr(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return an dep's category and package.
func (self *Cpv) Cpn() string {
	s := C.pkgcraft_dep_cpn(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an dep's category, package, and version.
func (self *Cpv) CPV() string {
	s := C.pkgcraft_dep_cpv(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Cpv) String() string {
	s := C.pkgcraft_dep_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Cpv) Hash() uint64 {
	if self._hash == 0 {
		self._hash = uint64(C.pkgcraft_dep_hash(self.ptr))
	}
	return self._hash
}

// Compare two deps returning -1, 0, or 1 if the first is less than, equal to,
// or greater than the second, respectively.
func (self *Cpv) Cmp(other depPtr) int {
	return int(C.pkgcraft_dep_cmp(self.ptr, other.p()))
}

// Determine if two deps intersect.
func (self *Cpv) Intersects(other depPtr) bool {
	return bool(C.pkgcraft_dep_intersects(self.ptr, other.p()))
}
