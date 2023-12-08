package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type DependencyUnit int

const (
	DependencyUnitDep DependencyUnit = iota
	DependencyUnitString
	DependencyUnitUri
)

type Dependency struct {
	ptr *C.Dependency
}

func depFromPtr(ptr *C.Dependency) *Dependency {
	obj := &Dependency{ptr}
	runtime.SetFinalizer(obj, func(self *Dependency) { C.pkgcraft_dependency_free(self.ptr) })
	return obj
}

func (self *Dependency) String() string {
	s := C.pkgcraft_dependency_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

type DependencySet struct {
	ptr *C.DependencySet
}

func dependencySetFromPtr(ptr *C.DependencySet) *DependencySet {
	obj := &DependencySet{ptr}
	runtime.SetFinalizer(obj, func(self *DependencySet) { C.pkgcraft_dependency_set_free(self.ptr) })
	return obj
}

func (self *DependencySet) String() string {
	s := C.pkgcraft_dependency_set_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
