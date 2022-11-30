package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/hashicorp/golang-lru"
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

var atom_cache, _ = lru.New(10000)

type Atom struct {
	ptr *C.Atom
	// cached fields
	_category string
	_package  string
	_version  *Version
	_hash     uint64
}

type Pair[T, U any] struct {
	First  T
	Second U
}

func newAtom(s string, eapi *Eapi) (*Atom, error) {
	var eapi_ptr *C.Eapi
	if eapi == nil {
		eapi_ptr = nil
	} else {
		eapi_ptr = eapi.ptr
	}

	atom_str := C.CString(s)
	defer C.free(unsafe.Pointer(atom_str))
	ptr := C.pkgcraft_atom_new(atom_str, eapi_ptr)

	if ptr != nil {
		atom := &Atom{ptr: ptr}
		runtime.SetFinalizer(atom, func(a *Atom) { C.pkgcraft_atom_free(a.ptr) })
		return atom, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

// Parse a string into an atom using the latest EAPI.
func NewAtom(s string) (*Atom, error) {
	return newAtom(s, nil)
}

// Parse a string into an atom using a specific EAPI.
func NewAtomWithEapi(s string, eapi *Eapi) (*Atom, error) {
	return newAtom(s, eapi)
}

func newCachedAtom(s string, eapi *Eapi) (*Atom, error) {
	key := Pair[string, *Eapi]{s, eapi}
	if v, ok := atom_cache.Get(key); ok {
		return v.(*Atom), nil
	} else {
		atom, err := newAtom(s, eapi)
		if err == nil {
			atom_cache.Add(key, atom)
		}
		return atom, err
	}
}

// Return a cached Atom if one exists, otherwise return a new instance.
func NewAtomCached(s string) (*Atom, error) {
	return newCachedAtom(s, nil)
}

// Return a cached Atom if one exists, otherwise parse using a specific EAPI.
func NewAtomCachedWithEapi(s string, eapi *Eapi) (*Atom, error) {
	return newCachedAtom(s, eapi)
}

// Return an atom's category.
func (a *Atom) Category() string {
	if a._category == "" {
		s := C.pkgcraft_atom_category(a.ptr)
		defer C.pkgcraft_str_free(s)
		a._category = C.GoString(s)
	}
	return a._category
}

// Return an atom's package name.
func (a *Atom) PN() string {
	if a._package == "" {
		s := C.pkgcraft_atom_package(a.ptr)
		defer C.pkgcraft_str_free(s)
		a._package = C.GoString(s)
	}
	return a._package
}

// Return an atom's version.
func (a *Atom) Version() *Version {
	if a._version == nil {
		ptr := C.pkgcraft_atom_version(a.ptr)
		if ptr != nil {
			a._version = &Version{ptr}
		} else {
			a._version = &Version{}
		}
	}
	return a._version
}

// Return an atom's revision.
func (a *Atom) Revision() string {
	s := C.pkgcraft_atom_revision(a.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's blocker.
func (a *Atom) Blocker() Blocker {
	i := C.pkgcraft_atom_blocker(a.ptr)
	return Blocker(i)
}

// Return an atom's slot.
func (a *Atom) Slot() string {
	s := C.pkgcraft_atom_slot(a.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's subslot.
func (a *Atom) Subslot() string {
	s := C.pkgcraft_atom_subslot(a.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's slot operator.
func (a *Atom) SlotOp() SlotOperator {
	i := C.pkgcraft_atom_slot_op(a.ptr)
	return SlotOperator(i)
}

// Return an atom's USE deps.
func (a *Atom) UseDeps() []string {
	var length C.size_t
	array := C.pkgcraft_atom_use_deps(a.ptr, &length)
	use_slice := unsafe.Slice(array, length)
	use := []string{}
	for _, s := range use_slice {
		use = append(use, C.GoString(s))
	}
	defer C.pkgcraft_str_array_free(array, length)
	return use
}

// Return an atom's repo.
func (a *Atom) Repo() string {
	s := C.pkgcraft_atom_repo(a.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return the concatenated string of an atom's category and package.
func (a *Atom) Key() string {
	s := C.pkgcraft_atom_key(a.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return the concatenated string of an atom's category, package, and version.
func (a *Atom) CPV() string {
	s := C.pkgcraft_atom_cpv(a.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Compare an atom with another atom returning -1, 0, or 1 if the first atom is
// less than, equal to, or greater than the second atom, respectively.
func (a1 *Atom) Cmp(a2 *Atom) int {
	return int(C.pkgcraft_atom_cmp(a1.ptr, a2.ptr))
}

type Atoms []*Atom

func (s Atoms) Len() int           { return len(s) }
func (s Atoms) Less(i, j int) bool { return s[i].Cmp(s[j]) == -1 }
func (s Atoms) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (a *Atom) String() string {
	s := C.pkgcraft_atom_str(a.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (a *Atom) Hash() uint64 {
	if a._hash == 0 {
		a._hash = uint64(C.pkgcraft_atom_hash(a.ptr))
	}
	return a._hash
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
		cpv := &Cpv{Atom{ptr: ptr}}
		runtime.SetFinalizer(cpv, func(cpv *Cpv) { C.pkgcraft_atom_free(cpv.ptr) })
		return cpv, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}
