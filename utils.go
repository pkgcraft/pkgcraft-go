package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"unsafe"
)

// Convert a slice of Go strings to an array of C strings.
func sliceToCharArray(vals []string) (**C.char, C.size_t) {
	c_strs := C.malloc(C.size_t(len(vals)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	a := (*[1<<30 - 1]*C.char)(c_strs)
	for i, s := range vals {
		a[i] = C.CString(s)
	}
	return (**C.char)(c_strs), C.size_t(len(vals))
}
