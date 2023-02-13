package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type DepSpecUnit int

const (
	DepSpecUnitDep DepSpecUnit = iota
	DepSpecUnitString
	DepSpecUnitUri
)

type DepSpec struct {
	ptr  *C.DepSpec
}

func depFromPtr(ptr *C.DepSpec) *DepSpec {
	obj := &DepSpec{ptr}
	runtime.SetFinalizer(obj, func(self *DepSpec) { C.pkgcraft_dep_spec_free(self.ptr) })
	return obj
}

func (self *DepSpec) String() string {
	s := C.pkgcraft_dep_spec_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
