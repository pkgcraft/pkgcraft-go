package atom

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type Blocker int

const (
	BlockerNone Blocker = iota - 1
	BlockerStrong
	BlockerWeak
)

type SlotOperator int

const (
	SlotOpNone SlotOperator = iota - 1
	SlotOpEqual
	SlotOpStar
)

// sentinel value for atoms with uncached version fields
var uncached_ver, _ = NewVersion("0")

type Atom struct {
	atom *C.Atom
	// cached fields
	_category string
	_package string
	_version *Version
}

func new_atom(s string, eapi string) (*Atom, error) {
	atom_str := C.CString(s)
	defer C.free(unsafe.Pointer(atom_str))

	var eapi_str *C.char
	if eapi == "" {
		eapi_str = nil
	} else {
		eapi_str = C.CString(eapi)
		defer C.free(unsafe.Pointer(eapi_str))
	}

	ptr := C.pkgcraft_atom_new(atom_str, eapi_str)

	if ptr != nil {
		atom := &Atom{atom: ptr, _version: uncached_ver}
		runtime.SetFinalizer(atom, func(a *Atom) { C.pkgcraft_atom_free(a.atom) })
		return atom, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
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
	if a._category == "" {
		s := C.pkgcraft_atom_category(a.atom)
		defer C.pkgcraft_str_free(s)
		a._category = C.GoString(s)
	}
	return a._category
}

// Return an atom's package name.
func (a *Atom) pn() string {
	if a._package == "" {
		s := C.pkgcraft_atom_package(a.atom)
		defer C.pkgcraft_str_free(s)
		a._package = C.GoString(s)
	}
	return a._package
}

// Return an atom's version.
func (a *Atom) version() *Version {
	if a._version == uncached_ver {
		ptr := C.pkgcraft_atom_version(a.atom)
		var ver *Version
		if ptr != nil {
			a._version = &Version{ptr}
		} else {
			a._version = ver
		}
	}
	return a._version
}

// Return an atom's revision.
func (a *Atom) revision() string {
	s := C.pkgcraft_atom_revision(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's blocker.
func (a *Atom) blocker() Blocker {
	i := C.pkgcraft_atom_blocker(a.atom)
	return Blocker(i)
}

// Return an atom's slot.
func (a *Atom) slot() string {
	s := C.pkgcraft_atom_slot(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's subslot.
func (a *Atom) subslot() string {
	s := C.pkgcraft_atom_subslot(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's slot operator.
func (a *Atom) slot_op() SlotOperator {
	i := C.pkgcraft_atom_slot_op(a.atom)
	return SlotOperator(i)
}

// Return an atom's USE deps.
func (a *Atom) use_deps() []string {
	var length C.size_t
	array := C.pkgcraft_atom_use_deps(a.atom, &length)
	use_slice := unsafe.Slice(array, length)
	use := []string{}
	for _, s := range use_slice {
		use = append(use, C.GoString(s))
	}
	defer C.pkgcraft_str_array_free(array, length)
	return use
}

// Return an atom's repo.
func (a *Atom) repo() string {
	s := C.pkgcraft_atom_repo(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return the concatenated string of an atom's category and package.
func (a *Atom) key() string {
	s := C.pkgcraft_atom_key(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return the concatenated string of an atom's category, package, and version.
func (a *Atom) cpv() string {
	s := C.pkgcraft_atom_cpv(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Compare an atom with another atom returning -1, 0, or 1 if the first atom is
// less than, equal to, or greater than the second atom, respectively.
func (a *Atom) cmp(b *Atom) int {
	return int(C.pkgcraft_atom_cmp(a.atom, b.atom))
}

type Atoms []*Atom

func (s Atoms) Len() int           { return len(s) }
func (s Atoms) Less(i, j int) bool { return s[i].cmp(s[j]) == -1 }
func (s Atoms) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (a *Atom) String() string {
	s := C.pkgcraft_atom_str(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (a *Atom) hash() uint64 {
	return uint64(C.pkgcraft_atom_hash(a.atom))
}

type Cpv struct {
    Atom
}

// Parse a CPV string into an atom.
func NewCpv(s string) (*Cpv, error) {
	cpv_str := C.CString(s)
	defer C.free(unsafe.Pointer(cpv_str))
	ptr := C.pkgcraft_cpv_new(cpv_str)

	if ptr != nil {
		cpv := &Cpv{Atom{atom: ptr, _version: uncached_ver}}
		runtime.SetFinalizer(cpv, func(cpv *Cpv) { C.pkgcraft_atom_free(cpv.atom) })
		return cpv, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}
