package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/hashicorp/golang-lru"
)

type Atom struct {
	Cpv
}

type Blocker int

const (
	BlockerNone Blocker = iota - 1
	BlockerStrong
	BlockerWeak
)

func BlockerFromString(s string) (Blocker, error) {
	c_str := C.CString(s)
	i := C.pkgcraft_atom_blocker_from_str(c_str)
	C.free(unsafe.Pointer(c_str))
	if i >= 0 {
		return Blocker(i), nil
	} else {
		return BlockerNone, fmt.Errorf("invalid blocker: %s", s)
	}
}

type SlotOperator int

const (
	SlotOpNone SlotOperator = iota - 1
	SlotOpEqual
	SlotOpStar
)

func SlotOperatorFromString(s string) (SlotOperator, error) {
	c_str := C.CString(s)
	i := C.pkgcraft_atom_slot_op_from_str(c_str)
	C.free(unsafe.Pointer(c_str))
	if i >= 0 {
		return SlotOperator(i), nil
	} else {
		return SlotOpNone, fmt.Errorf("invalid slot operator: %s", s)
	}
}

var atom_cache, _ = lru.New(10000)

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

	c_str := C.CString(s)
	ptr := C.pkgcraft_atom_new(c_str, eapi_ptr)
	C.free(unsafe.Pointer(c_str))

	if ptr != nil {
		atom := &Atom{Cpv{ptr: ptr}}
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
func (a *Atom) Use() []string {
	var length C.size_t
	array := C.pkgcraft_atom_use_deps(a.ptr, &length)
	use_slice := unsafe.Slice(array, length)
	var use []string
	for _, s := range use_slice {
		use = append(use, C.GoString(s))
	}
	C.pkgcraft_str_array_free(array, length)
	return use
}

// Return an atom's repo.
func (a *Atom) Repo() string {
	s := C.pkgcraft_atom_repo(a.ptr)
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
