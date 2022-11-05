package eapi

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"strconv"
	"unsafe"
)

var EapisOfficial = get_official_eapis()
var Eapis = get_eapis()
var EapiLatest = EapisOfficial[strconv.Itoa(len(EapisOfficial) - 1)]

// Convert an array of Eapi pointers to an (id, Eapi) mapping.
func eapis_to_map(eapis []*C.Eapi) map[string]Eapi {
	m := make(map[string]Eapi)
	for _, eapi := range eapis {
		s := C.pkgcraft_eapi_as_str(eapi)
		id := C.GoString(s)
		defer C.pkgcraft_str_free(s)
		m[id] = Eapi{eapi: eapi, _id: id}
	}
	return m
}

// Return the mapping of all official EAPIs.
func get_official_eapis() map[string]Eapi {
	var length C.size_t
	eapis := C.pkgcraft_eapis_official(&length)
	m := eapis_to_map(unsafe.Slice(eapis, length))
	defer C.pkgcraft_eapis_free(eapis, length)
	return m
}

// Return the mapping of all known EAPIs.
func get_eapis() map[string]Eapi {
	var length C.size_t
	eapis := C.pkgcraft_eapis(&length)
	m := eapis_to_map(unsafe.Slice(eapis, length))
	defer C.pkgcraft_eapis_free(eapis, length)
	return m
}

type Eapi struct {
	eapi *C.Eapi
	// cached fields
	_id string
}

// Return the string for an EAPI.
func (e *Eapi) String() string {
	return e._id
}

// Check if an EAPI has a given feature.
func (e *Eapi) has(s string) bool {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	return C.pkgcraft_eapi_has(e.eapi, cstr) == true
}
