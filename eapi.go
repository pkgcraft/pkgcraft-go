package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"strconv"
	"unsafe"
)

var EAPIS_OFFICIAL = getOfficialEapis()
var EAPIS = getEapis()
var EAPI_LATEST = EAPIS_OFFICIAL[strconv.Itoa(len(EAPIS_OFFICIAL)-1)]

// Convert an array of Eapi pointers to a mapping.
func eapisToMap(eapis []*C.Eapi, start int) map[string]*Eapi {
	m := make(map[string]*Eapi)
	for i, ptr := range eapis {
		if i >= start {
			s := C.pkgcraft_eapi_as_str(ptr)
			id := C.GoString(s)
			defer C.pkgcraft_str_free(s)
			m[id] = &Eapi{ptr, id}
		}
	}
	return m
}

// Return the mapping of all official EAPIs.
func getOfficialEapis() map[string]*Eapi {
	var length C.size_t
	eapis := C.pkgcraft_eapis_official(&length)
	m := eapisToMap(unsafe.Slice(eapis, length), 0)
	defer C.pkgcraft_eapis_free(eapis, length)
	return m
}

// Return the mapping of all known EAPIs.
func getEapis() map[string]*Eapi {
	var length C.size_t
	eapis := C.pkgcraft_eapis(&length)
	m := make(map[string]*Eapi)
	for k, v := range EAPIS_OFFICIAL {
		m[k] = v
	}
	for k, v := range eapisToMap(unsafe.Slice(eapis, length), len(m)) {
		m[k] = v
	}
	defer C.pkgcraft_eapis_free(eapis, length)
	return m
}

// Convert an EAPI range into an array of Eapi objects.
func EapiRange(s string) ([]*Eapi, error) {
	var length C.size_t
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	c_eapis := C.pkgcraft_eapis_range(cstr, &length)
	if c_eapis == nil {
		err := C.pkgcraft_error_last()
		defer C.pkgcraft_error_free(err)
		return nil, errors.New(C.GoString(err.message))
	}

	var eapis []*Eapi
	for _, ptr := range unsafe.Slice(c_eapis, length) {
		s := C.pkgcraft_eapi_as_str(ptr)
		id := C.GoString(s)
		defer C.pkgcraft_str_free(s)
		eapis = append(eapis, EAPIS[id])
	}
	return eapis, nil
}

type Eapi struct {
	ptr *C.Eapi
	// cached fields
	id string
}

// Return the string for an EAPI.
func (self *Eapi) String() string {
	return self.id
}

// Check if an EAPI has a given feature.
func (self *Eapi) Has(s string) bool {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	return bool(C.pkgcraft_eapi_has(self.ptr, cstr))
}

// Compare an Eapi with another Eapi chronologically returning -1, 0, or 1 if
// the first is less than, equal to, or greater than the second, respectively.
func (self *Eapi) Cmp(other *Eapi) int {
	return int(C.pkgcraft_eapi_cmp(self.ptr, other.ptr))
}
