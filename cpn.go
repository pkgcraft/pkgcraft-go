package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
	"unsafe"
)

type Cpn struct {
	ptr *C.Cpn
	// cached fields
	_category string
	_package  string
	_hash     uint64
}

type cpnPtr interface {
	p() *C.Cpn
}

func cpnFromPtr(ptr *C.Cpn) (*Cpn, error) {
	if ptr != nil {
		cpn := &Cpn{ptr: ptr}
		runtime.SetFinalizer(cpn, func(self *Cpn) { C.pkgcraft_cpn_free(self.ptr) })
		return cpn, nil
	} else {
		return nil, newPkgcraftError()
	}
}

// Parse a string into a Cpn object.
func NewCpn(s string) (*Cpn, error) {
	c_str := C.CString(s)
	defer C.free(unsafe.Pointer(c_str))
	ptr := C.pkgcraft_cpn_new(c_str)
	return cpnFromPtr(ptr)
}

func (self *Cpn) p() *C.Cpn {
	return self.ptr
}

// Return an Cpn's category.
func (self *Cpn) Category() string {
	if self._category == "" {
		s := C.pkgcraft_cpn_category(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._category = C.GoString(s)
	}
	return self._category
}

// Return a Cpn's package name.
func (self *Cpn) Package() string {
	if self._package == "" {
		s := C.pkgcraft_cpn_package(self.ptr)
		defer C.pkgcraft_str_free(s)
		self._package = C.GoString(s)
	}
	return self._package
}

func (self *Cpn) String() string {
	s := C.pkgcraft_cpn_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

func (self *Cpn) Hash() uint64 {
	if self._hash == 0 {
		self._hash = uint64(C.pkgcraft_cpn_hash(self.ptr))
	}
	return self._hash
}

// Compare two Cpns returning -1, 0, or 1 if the first is less than, equal to,
// or greater than the second, respectively.
func (self *Cpn) Cmp(other *Cpn) int {
	return int(C.pkgcraft_cpn_cmp(self.ptr, other.ptr))
}
