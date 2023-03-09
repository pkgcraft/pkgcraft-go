package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"unsafe"
)

var EAPIS_OFFICIAL = getOfficialEapis()
var EAPIS = getEapis()
var EAPI_LATEST_OFFICIAL *Eapi
var EAPI_LATEST *Eapi

// Convert an array of Eapi pointers to a slice of Eapi objects.
func eapisToSlice(c_eapis []*C.Eapi, start int) []*Eapi {
	var eapis []*Eapi
	for i, ptr := range c_eapis {
		if i >= start {
			s := C.pkgcraft_eapi_as_str(ptr)
			id := C.GoString(s)
			defer C.pkgcraft_str_free(s)
			eapis = append(eapis, &Eapi{ptr, id})
		}
	}
	return eapis
}

// Return the mapping of all official EAPIs.
func getOfficialEapis() map[string]*Eapi {
	var length C.size_t
	c_eapis := C.pkgcraft_eapis_official(&length)
	eapis := eapisToSlice(unsafe.Slice(c_eapis, length), 0)
	defer C.pkgcraft_eapis_free(c_eapis, length)

	// set global alias for the most recent, official EAPI
	EAPI_LATEST_OFFICIAL = eapis[len(eapis)-1]

	m := make(map[string]*Eapi)
	for _, eapi := range eapis {
		m[eapi.id] = eapi
	}

	return m
}

// Return the mapping of all known EAPIs.
func getEapis() map[string]*Eapi {
	var eapis []*Eapi
	var length C.size_t
	c_eapis := C.pkgcraft_eapis(&length)

	// copy official Eapi objects
	for _, eapi := range EAPIS_OFFICIAL {
		eapis = append(eapis, eapi)
	}

	// append unofficial Eapi objects
	unofficial_eapis := eapisToSlice(unsafe.Slice(c_eapis, length), len(eapis))
	eapis = append(eapis, unofficial_eapis...)
	defer C.pkgcraft_eapis_free(c_eapis, length)

	// set global alias for the most recent EAPI
	EAPI_LATEST = unofficial_eapis[len(unofficial_eapis)-1]

	m := make(map[string]*Eapi)
	for _, eapi := range eapis {
		m[eapi.id] = eapi
	}

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
