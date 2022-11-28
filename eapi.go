package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"strconv"
	"unsafe"
)

var EAPIS_OFFICIAL = get_official_eapis()
var EAPIS = get_eapis()
var EAPI_LATEST = EAPIS_OFFICIAL[strconv.Itoa(len(EAPIS_OFFICIAL)-1)]

// Convert an array of Eapi pointers to a mapping.
func eapis_to_map(eapis []*C.Eapi, start int) map[string]*Eapi {
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
func get_official_eapis() map[string]*Eapi {
	var length C.size_t
	eapis := C.pkgcraft_eapis_official(&length)
	m := eapis_to_map(unsafe.Slice(eapis, length), 0)
	defer C.pkgcraft_eapis_free(eapis, length)
	return m
}

// Return the mapping of all known EAPIs.
func get_eapis() map[string]*Eapi {
	var length C.size_t
	eapis := C.pkgcraft_eapis(&length)
	m := make(map[string]*Eapi)
	for k, v := range EAPIS_OFFICIAL {
		m[k] = v
	}
	for k, v := range eapis_to_map(unsafe.Slice(eapis, length), len(m)) {
		m[k] = v
	}
	defer C.pkgcraft_eapis_free(eapis, length)
	return m
}

type Eapi struct {
	ptr *C.Eapi
	// cached fields
	id string
}

// Return the string for an EAPI.
func (e *Eapi) String() string {
	return e.id
}

// Check if an EAPI has a given feature.
func (e *Eapi) Has(s string) bool {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	return C.pkgcraft_eapi_has(e.ptr, cstr) == true
}
