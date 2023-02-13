package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"unsafe"

	"github.com/hashicorp/golang-lru/v2"
)

type Dep struct {
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

var dep_cache, _ = lru.New[Pair[string, *Eapi], *Dep](10000)

func newDep(s string, eapi *Eapi) (*Dep, error) {
	var eapi_ptr *C.Eapi
	if eapi == nil {
		eapi_ptr = nil
	} else {
		eapi_ptr = eapi.ptr
	}

	c_str := C.CString(s)
	ptr := C.pkgcraft_dep_new(c_str, eapi_ptr)
	C.free(unsafe.Pointer(c_str))

	cpv, err := cpvFromPtr(ptr)
	if cpv != nil {
		return &Dep{cpv}, nil
	} else {
		return nil, err
	}
}

// Parse a string into a Dep using the latest EAPI.
func NewDep(s string) (*Dep, error) {
	return newDep(s, nil)
}

// Parse a string into a Dep using a specific EAPI.
func NewDepWithEapi(s string, eapi *Eapi) (*Dep, error) {
	return newDep(s, eapi)
}

func newCachedDep(s string, eapi *Eapi) (*Dep, error) {
	key := Pair[string, *Eapi]{s, eapi}
	if dep, ok := dep_cache.Get(key); ok {
		return dep, nil
	} else {
		dep, err := newDep(s, eapi)
		if err == nil {
			dep_cache.Add(key, dep)
		}
		return dep, err
	}
}

// Return a cached Dep if one exists, otherwise return a new instance.
func NewDepCached(s string) (*Dep, error) {
	return newCachedDep(s, nil)
}

// Return a cached Dep if one exists, otherwise parse using a specific EAPI.
func NewDepCachedWithEapi(s string, eapi *Eapi) (*Dep, error) {
	return newCachedDep(s, eapi)
}

// Get the blocker of a package dependency.
func (self *Dep) Blocker() Blocker {
	i := C.pkgcraft_dep_blocker(self.ptr)
	return Blocker(i)
}

// Get the slot of a package dependency.
func (self *Dep) Slot() string {
	s := C.pkgcraft_dep_slot(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Get the subslot of a package dependency.
func (self *Dep) Subslot() string {
	s := C.pkgcraft_dep_subslot(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Get the slot operator of a package dependency.
func (self *Dep) SlotOp() SlotOperator {
	i := C.pkgcraft_dep_slot_op(self.ptr)
	return SlotOperator(i)
}

// Get the USE dependencies of a package dependency.
func (self *Dep) Use() []string {
	var length C.size_t
	array := C.pkgcraft_dep_use_deps(self.ptr, &length)
	use_slice := unsafe.Slice(array, length)
	var use []string
	for _, s := range use_slice {
		use = append(use, C.GoString(s))
	}
	C.pkgcraft_str_array_free(array, length)
	return use
}

// Get the repo of a package dependency.
func (self *Dep) Repo() string {
	s := C.pkgcraft_dep_repo(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Compare two package dependencies returning -1, 0, or 1 if the first is
// less than, equal to, or greater than the second, respectively.
func (self *Dep) Cmp(other *Dep) int {
	return int(C.pkgcraft_dep_cmp(self.ptr, other.ptr))
}
