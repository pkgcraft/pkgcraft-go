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

func atom_free(a *Atom) {
	C.pkgcraft_atom_free(a.atom)
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

	ptr := C.pkgcraft_atom(atom_str, eapi_str)

	if ptr != nil {
		atom := &Atom{ptr}
		runtime.SetFinalizer(atom, atom_free)
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
	s := C.pkgcraft_atom_category(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's package name.
func (a *Atom) pn() string {
	s := C.pkgcraft_atom_package(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's version.
func (a *Atom) version() string {
	s := C.pkgcraft_atom_version(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return an atom's revision.
func (a *Atom) revision() string {
	s := C.pkgcraft_atom_revision(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
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
func (a *Atom) slot_op() string {
	s := C.pkgcraft_atom_slot_op(a.atom)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
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

type Version struct {
	version *C.Version
}

func version_free(v *Version) {
	C.pkgcraft_version_free(v.version)
}

// Parse a string into a version.
func NewVersion(s string) (*Version, error) {
	ver_str := C.CString(s)
	defer C.free(unsafe.Pointer(ver_str))
	ptr := C.pkgcraft_version(ver_str)

	if ptr != nil {
		ver := &Version{ptr}
		runtime.SetFinalizer(ver, version_free)
		return ver, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

// Return a version's revision.
func (v *Version) revision() string {
	s := C.pkgcraft_version_revision(v.version)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Compare a version with another version returning -1, 0, or 1 if the first
// version is less than, equal to, or greater than the second version,
// respectively.
func (a *Version) cmp(b *Version) int {
	return int(C.pkgcraft_version_cmp(a.version, b.version))
}

type Versions []*Version

func (s Versions) Len() int           { return len(s) }
func (s Versions) Less(i, j int) bool { return s[i].cmp(s[j]) == -1 }
func (s Versions) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (v *Version) String() string {
	s := C.pkgcraft_version_str(v.version)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (v *Version) hash() uint64 {
	return uint64(C.pkgcraft_version_hash(v.version))
}
