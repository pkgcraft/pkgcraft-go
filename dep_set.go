package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type DepSet struct {
	ptr  *C.DepSet
}

func depSetFromPtr(ptr *C.DepSet) *DepSet {
	obj := &DepSet{ptr}
	runtime.SetFinalizer(obj, func(self *DepSet) { C.pkgcraft_dep_set_free(self.ptr) })
	return obj
}

func (self *DepSet) String() string {
	s := C.pkgcraft_dep_set_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
