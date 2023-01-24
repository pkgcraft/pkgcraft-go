package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"unsafe"

	"github.com/hashicorp/golang-lru/v2"
)

type Atom struct {
	*Cpv
}

type Blocker int

const (
	BlockerNone Blocker = iota
	BlockerStrong
	BlockerWeak
)

type SlotOperator int

const (
	SlotOpNone SlotOperator = iota
	SlotOpEqual
	SlotOpStar
)

type Pair[T, U any] struct {
	First  T
	Second U
}

var atom_cache, _ = lru.New[Pair[string, *Eapi], *Atom](10000)

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

	cpv, err := cpvFromPtr(ptr)
	if cpv != nil {
		return &Atom{cpv}, nil
	} else {
		return nil, err
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
	if atom, ok := atom_cache.Get(key); ok {
		return atom, nil
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
func (self *Atom) Blocker() Blocker {
	i := C.pkgcraft_atom_blocker(self.ptr)
	return Blocker(i)
}

// Return an atom's slot.
func (self *Atom) Slot() string {
	s := C.pkgcraft_atom_slot(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's subslot.
func (self *Atom) Subslot() string {
	s := C.pkgcraft_atom_subslot(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's slot operator.
func (self *Atom) SlotOp() SlotOperator {
	i := C.pkgcraft_atom_slot_op(self.ptr)
	return SlotOperator(i)
}

// Return an atom's USE deps.
func (self *Atom) Use() []string {
	var length C.size_t
	array := C.pkgcraft_atom_use_deps(self.ptr, &length)
	use_slice := unsafe.Slice(array, length)
	var use []string
	for _, s := range use_slice {
		use = append(use, C.GoString(s))
	}
	C.pkgcraft_str_array_free(array, length)
	return use
}

// Return an atom's repo.
func (self *Atom) Repo() string {
	s := C.pkgcraft_atom_repo(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return the concatenated string of an atom's category, package, and version.
func (self *Atom) CPV() string {
	s := C.pkgcraft_atom_cpv(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Compare an atom with another atom returning -1, 0, or 1 if the first atom is
// less than, equal to, or greater than the second atom, respectively.
func (self *Atom) Cmp(other *Atom) int {
	return int(C.pkgcraft_atom_cmp(self.ptr, other.ptr))
}
