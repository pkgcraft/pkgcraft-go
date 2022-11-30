package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

// Return the pkgcraft library version.
func pkgcraftLibVersion() string {
	s := C.pkgcraft_lib_version()
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
