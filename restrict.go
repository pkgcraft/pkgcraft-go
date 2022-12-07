package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type Restrict struct {
	ptr *C.Restrict
}

func restrictFromPtr(ptr *C.Restrict) *Restrict {
	restrict := &Restrict{ptr}
	runtime.SetFinalizer(restrict, func(self *Restrict) { C.pkgcraft_restrict_free(self.ptr) })
	return restrict
}

// Return a new restriction from a given object.
func NewRestrict(obj interface{}) (*Restrict, error) {
	ptr, err := objectToRestrict(obj)
	if ptr != nil {
		restrict := restrictFromPtr(ptr)
		return restrict, nil
	} else {
		return nil, err
	}
}

// Try to convert a string to a restriction.
func stringToRestrict(s string) (*C.Restrict, error) {
	if cpv, _ := NewCpv(s); cpv != nil {
		return C.pkgcraft_atom_restrict(cpv.ptr), nil
	} else if atom, _ := NewAtomCached(s); atom != nil {
		return C.pkgcraft_atom_restrict(atom.ptr), nil
	} else {
		c_str := C.CString(s)
		defer C.free(unsafe.Pointer(c_str))
		if ptr := C.pkgcraft_restrict_parse_dep(c_str); ptr != nil {
			return ptr, nil
		} else if ptr := C.pkgcraft_restrict_parse_pkg(c_str); ptr != nil {
			return ptr, nil
		}
	}
	return nil, fmt.Errorf("invalid restriction string: %s", s)
}

// Try to convert an object to a restriction.
func objectToRestrict(obj interface{}) (*C.Restrict, error) {
	switch obj := obj.(type) {
	case *Cpv:
		return C.pkgcraft_atom_restrict(obj.ptr), nil
	case *BasePkg:
		return C.pkgcraft_pkg_restrict(obj.ptr), nil
	case string:
		return stringToRestrict(obj)
	default:
		return nil, fmt.Errorf("unsupported restrict type: %t", obj)
	}
}

// Create a new restriction combining two restrictions via logical AND.
func (self *Restrict) And(other *Restrict) *Restrict {
	ptr := C.pkgcraft_restrict_and(self.ptr, other.ptr)
	return restrictFromPtr(ptr)
}

// Create a new restriction combining two restrictions via logical OR.
func (self *Restrict) Or(other *Restrict) *Restrict {
	ptr := C.pkgcraft_restrict_or(self.ptr, other.ptr)
	return restrictFromPtr(ptr)
}

// Create a new restriction combining two restrictions via logical XOR.
func (self *Restrict) Xor(other *Restrict) *Restrict {
	ptr := C.pkgcraft_restrict_xor(self.ptr, other.ptr)
	return restrictFromPtr(ptr)
}

// Create a new restriction inverting a restriction via logical NOT.
func (self *Restrict) Not() *Restrict {
	ptr := C.pkgcraft_restrict_not(self.ptr)
	return restrictFromPtr(ptr)
}
