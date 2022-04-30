package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type Atom struct {
	atom *C.Atom
}

func new_atom(s string, eapi string) (*Atom, error) {
	atom_str := C.CString(s)
	defer C.free(unsafe.Pointer(atom_str))

	var eapi_str *C.char
	if eapi != "" {
		eapi_str = C.CString(eapi)
		defer C.free(unsafe.Pointer(eapi_str))
	} else {
		eapi_str = nil
	}

	atom := &Atom{C.pkgcraft_atom(atom_str, eapi_str)}
	runtime.SetFinalizer(atom, func(a *Atom) { C.pkgcraft_atom_free(a.atom) })

	if atom.atom != nil {
		return atom, nil
	} else {
		return atom, errors.New(C.GoString(C.pkgcraft_last_error()))
	}
}

// Parse a string into an atom using the latest EAPI.
func NewAtom(s string) (*Atom, error) {
	return new_atom(s, "")
}

// Parse a string into an atom using a specific EAPI.
func NewAtomWithEapi(s string, eapi string) (*Atom, error) {
	return new_atom(s, eapi)
}

// Return an atom's category.
func (a *Atom) category() string {
	s := C.pkgcraft_atom_category(a.atom)
	defer C.pkgcraft_free_str(s)
	return C.GoString(s)
}

// Return an atom's package name.
func (a *Atom) pn() string {
	s := C.pkgcraft_atom_package(a.atom)
	defer C.pkgcraft_free_str(s)
	return C.GoString(s)
}

// Return an atom's version.
func (a *Atom) version() string {
	s := C.pkgcraft_atom_version(a.atom)
	defer C.pkgcraft_free_str(s)
	return C.GoString(s)
}

// Return an atom's slot.
func (a *Atom) slot() string {
	s := C.pkgcraft_atom_slot(a.atom)
	defer C.pkgcraft_free_str(s)
	return C.GoString(s)
}

// Return an atom's subslot.
func (a *Atom) subslot() string {
	s := C.pkgcraft_atom_subslot(a.atom)
	defer C.pkgcraft_free_str(s)
	return C.GoString(s)
}
