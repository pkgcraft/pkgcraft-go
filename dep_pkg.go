package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/hashicorp/golang-lru/v2"
)

type Dep struct {
	ptr *C.Dep
	// cached fields
	_category string
	_package  string
	_version  *Version
	_hash     uint64
}

type Blocker int

const (
	BlockerNone Blocker = iota
	BlockerStrong
	BlockerWeak
)

func BlockerFromString(s string) (Blocker, error) {
	c_str := C.CString(s)
	i := C.pkgcraft_dep_blocker_from_str(c_str)
	C.free(unsafe.Pointer(c_str))
	if i > 0 {
		return Blocker(i), nil
	} else {
		return BlockerNone, fmt.Errorf("invalid blocker: %s", s)
	}
}

type SlotOperator int

const (
	SlotOpNone SlotOperator = iota
	SlotOpEqual
	SlotOpStar
)

func SlotOperatorFromString(s string) (SlotOperator, error) {
	c_str := C.CString(s)
	i := C.pkgcraft_dep_slot_op_from_str(c_str)
	C.free(unsafe.Pointer(c_str))
	if i > 0 {
		return SlotOperator(i), nil
	} else {
		return SlotOpNone, fmt.Errorf("invalid slot operator: %s", s)
	}
}

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

	if ptr != nil {
		dep := &Dep{ptr: ptr}
		runtime.SetFinalizer(dep, func(self *Dep) { C.pkgcraft_dep_free(self.ptr) })
		return dep, nil
	} else {
		return nil, newPkgcraftError()
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

// Return a package dependency's category.
func (self *Dep) Category() string {
	if self._category == "" {
		s := C.pkgcraft_dep_category(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._category = C.GoString(s)
	}
	return self._category
}

// Return a package dependency's package.
func (self *Dep) Package() string {
	if self._package == "" {
		s := C.pkgcraft_dep_package(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._package = C.GoString(s)
	}
	return self._package
}

// Return a package dependency's version.
func (self *Dep) Version() *Version {
	if self._version == nil {
		ptr := C.pkgcraft_dep_version(self.ptr)
		if ptr != nil {
			self._version, _ = versionFromPtr(ptr)
		} else {
			self._version = &Version{}
		}
	}
	return self._version
}

// Return a package dependency's revision.
func (self *Dep) Revision() *Revision {
	version := self.Version()
	if *version != (Version{}) {
		return version.Revision()
	}
	return &Revision{}
}

// Return a package dependency's slot.
func (self *Dep) Slot() string {
	s := C.pkgcraft_dep_slot(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package dependency's subslot.
func (self *Dep) Subslot() string {
	s := C.pkgcraft_dep_subslot(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package dependency's slot operator.
func (self *Dep) SlotOp() SlotOperator {
	i := C.pkgcraft_dep_slot_op(self.ptr)
	return SlotOperator(i)
}

// Return a package dependency's USE flag dependencies.
func (self *Dep) Use() []string {
	var length C.size_t
	array := C.pkgcraft_dep_use_deps_str(self.ptr, &length)
	use_slice := unsafe.Slice(array, length)
	var use []string
	for _, s := range use_slice {
		use = append(use, C.GoString(s))
	}
	C.pkgcraft_str_array_free(array, length)
	return use
}

// Return a package dependency's repository.
func (self *Dep) Repo() string {
	s := C.pkgcraft_dep_repo(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package dependency's package and version.
func (self *Dep) P() string {
	c_str := C.pkgcraft_dep_p(self.ptr)
	defer C.pkgcraft_str_free(c_str)
	return C.GoString(c_str)
}

// Return a package dependency's package, version, and revision.
func (self *Dep) Pf() string {
	c_str := C.pkgcraft_dep_pf(self.ptr)
	defer C.pkgcraft_str_free(c_str)
	return C.GoString(c_str)
}

// Return a package dependency's revision.
func (self *Dep) Pr() string {
	c_str := C.pkgcraft_dep_pr(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return a package dependency's version.
func (self *Dep) Pv() string {
	c_str := C.pkgcraft_dep_pv(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return a package dependency's version and revision.
func (self *Dep) Pvr() string {
	c_str := C.pkgcraft_dep_pvr(self.ptr)
	if c_str != nil {
		defer C.pkgcraft_str_free(c_str)
		return C.GoString(c_str)
	}
	return ""
}

// Return a package dependency's Cpn.
func (self *Dep) Cpn() *Cpn {
	ptr := C.pkgcraft_dep_cpn(self.ptr)
	cpn, _ := cpnFromPtr(ptr)
	return cpn
}

// Return the Cpv of a package dependency if one exists.
func (self *Dep) Cpv() *Cpv {
	ptr := C.pkgcraft_dep_cpv(self.ptr)
	if ptr != nil {
		cpv, _ := cpvFromPtr(ptr)
		return cpv
	}
	return nil
}

func (self *Dep) String() string {
	s := C.pkgcraft_dep_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Dep) Hash() uint64 {
	if self._hash == 0 {
		self._hash = uint64(C.pkgcraft_dep_hash(self.ptr))
	}
	return self._hash
}

// Compare two package dependencies returning -1, 0, or 1 if the first is
// less than, equal to, or greater than the second, respectively.
func (self *Dep) Cmp(other *Dep) int {
	return int(C.pkgcraft_dep_cmp(self.ptr, other.ptr))
}

// Determine if two Cpv or Dep objects intersect.
func (self *Dep) Intersects(other interface{}) bool {
	switch other := other.(type) {
	case *Cpv:
		return bool(C.pkgcraft_dep_intersects_cpv(self.ptr, other.ptr))
	case *Dep:
		return bool(C.pkgcraft_dep_intersects(self.ptr, other.ptr))
	default:
		return false
	}
}
