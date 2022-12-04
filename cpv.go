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

// Parse a CPV string into an atom.
func NewCpv(s string) (*Cpv, error) {
	c_str := C.CString(s)
	defer C.free(unsafe.Pointer(c_str))
	ptr := C.pkgcraft_cpv_new(c_str)

	if ptr != nil {
		cpv := &Cpv{ptr: ptr}
		runtime.SetFinalizer(cpv, func(cpv *Cpv) { C.pkgcraft_atom_free(cpv.ptr) })
		return cpv, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

// Return an atom's category.
func (cpv *Cpv) Category() string {
	if cpv._category == "" {
		s := C.pkgcraft_atom_category(cpv.ptr)
		defer C.pkgcraft_str_free(s)
		cpv._category = C.GoString(s)
	}
	return cpv._category
}

// Return an atom's package name.
func (cpv *Cpv) Package() string {
	if cpv._package == "" {
		s := C.pkgcraft_atom_package(cpv.ptr)
		defer C.pkgcraft_str_free(s)
		cpv._package = C.GoString(s)
	}
	return cpv._package
}

// Return an atom's version.
func (cpv *Cpv) Version() *Version {
	if cpv._version == nil {
		ptr := C.pkgcraft_atom_version(cpv.ptr)
		if ptr != nil {
			cpv._version = &Version{ptr}
		} else {
			cpv._version = &Version{}
		}
	}
	return cpv._version
}

// Return an atom's revision.
func (cpv *Cpv) Revision() string {
	s := C.pkgcraft_atom_revision(cpv.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return the concatenated string of an atom's category and package.
func (cpv *Cpv) Key() string {
	s := C.pkgcraft_atom_key(cpv.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (cpv *Cpv) String() string {
	s := C.pkgcraft_atom_str(cpv.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (cpv *Cpv) Hash() uint64 {
	if cpv._hash == 0 {
		cpv._hash = uint64(C.pkgcraft_atom_hash(cpv.ptr))
	}
	return cpv._hash
}

// Compare an atom with another atom returning -1, 0, or 1 if the first atom is
// less than, equal to, or greater than the second atom, respectively.
func (c1 *Cpv) Cmp(c2 *Cpv) int {
	return int(C.pkgcraft_atom_cmp(c1.ptr, c2.ptr))
}
