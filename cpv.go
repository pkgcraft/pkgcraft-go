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
	ptr *C.Atom
	// cached fields
	_category string
	_package  string
	_version  *Version
	_hash     uint64
}

func cpvFromPtr(ptr *C.Atom) (*Cpv, error) {
	if ptr != nil {
		cpv := &Cpv{ptr: ptr}
		runtime.SetFinalizer(cpv, func(self *Cpv) { C.pkgcraft_atom_free(self.ptr) })
		return cpv, nil
	} else {
		err := C.pkgcraft_error_last()
		defer C.pkgcraft_error_free(err)
		return nil, errors.New(C.GoString(err.message))
	}
}

// Parse a CPV string into an atom.
func NewCpv(s string) (*Cpv, error) {
	c_str := C.CString(s)
	defer C.free(unsafe.Pointer(c_str))
	ptr := C.pkgcraft_cpv_new(c_str)
	return cpvFromPtr(ptr)
}

// Return an atom's category.
func (self *Cpv) Category() string {
	if self._category == "" {
		s := C.pkgcraft_atom_category(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._category = C.GoString(s)
	}
	return self._category
}

// Return an atom's package name.
func (self *Cpv) Package() string {
	if self._package == "" {
		s := C.pkgcraft_atom_package(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._package = C.GoString(s)
	}
	return self._package
}

// Return an atom's version.
func (self *Cpv) Version() *Version {
	if self._version == nil {
		ptr := C.pkgcraft_atom_version(self.ptr)
		if ptr != nil {
			self._version, _ = versionFromPtr(ptr)
		} else {
			self._version = &Version{}
		}
	}
	return self._version
}

// Return an atom's revision.
func (self *Cpv) Revision() string {
	s := C.pkgcraft_atom_revision(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return the concatenated string of an atom's category and package.
func (self *Cpv) Cpn() string {
	s := C.pkgcraft_atom_cpn(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Cpv) String() string {
	s := C.pkgcraft_atom_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Cpv) Hash() uint64 {
	if self._hash == 0 {
		self._hash = uint64(C.pkgcraft_atom_hash(self.ptr))
	}
	return self._hash
}

// Compare an atom with another atom returning -1, 0, or 1 if the first atom is
// less than, equal to, or greater than the second atom, respectively.
func (self *Cpv) Cmp(other *Cpv) int {
	return int(C.pkgcraft_atom_cmp(self.ptr, other.ptr))
}
