package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

type Restrict struct {
	ptr *C.Restrict
}

// Return a new config for the system.
func NewRestrict(obj interface{}) (*Restrict, error) {
	ptr, err := objectToRestrict(obj)
	if ptr != nil {
		restrict := &Restrict{ptr}
		runtime.SetFinalizer(restrict, func(r *Restrict) { C.pkgcraft_restrict_free(r.ptr) })
		return restrict, nil
	} else {
		return nil, err
	}
}

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
		} else {
			s := C.pkgcraft_last_error()
			defer C.pkgcraft_str_free(s)
			return nil, errors.New(C.GoString(s))
		}
	}
	return nil, fmt.Errorf("invalid restriction string: %s", s)
}

func objectToRestrict(obj interface{}) (*C.Restrict, error) {
	switch obj.(type) {
	case *Atom:
		return C.pkgcraft_atom_restrict(obj.(*Atom).ptr), nil
	case string:
		return stringToRestrict(obj.(string))
	default:
		return nil, fmt.Errorf("unsupported restrict type: %t", obj)
	}
}
