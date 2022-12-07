package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type DepSetKind int

const (
	DepSetAtom DepSetKind = iota
	DepSetString
	DepSetUri
)

type DepRestrict struct {
	ptr  *C.DepRestrict
	kind DepSetKind
}

func depRestrictFromPtr(ptr *C.DepRestrict, kind DepSetKind) *DepRestrict {
	obj := &DepRestrict{ptr, kind}
	runtime.SetFinalizer(obj, func(o *DepRestrict) { C.pkgcraft_deprestrict_free(o.ptr) })
	return obj
}

func (self *DepRestrict) String() string {
	s := C.pkgcraft_deprestrict_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

type DepSet struct {
	ptr  *C.DepSet
	kind DepSetKind
}

func depSetFromPtr(ptr *C.DepSet, kind DepSetKind) *DepSet {
	obj := &DepSet{ptr, kind}
	runtime.SetFinalizer(obj, func(o *DepSet) { C.pkgcraft_depset_free(o.ptr) })
	return obj
}

func (self *DepSet) String() string {
	s := C.pkgcraft_depset_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
