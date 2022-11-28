package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type EbuildPkg struct {
	*BasePkg
}

// Return a package's path.
func (p *EbuildPkg) Path() string {
	s := C.pkgcraft_pkg_ebuild_path(p.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
