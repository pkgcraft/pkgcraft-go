package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type BasePkg struct {
	ptr    *C.Pkg
	format PkgFormat
}

// Return a package's atom.
func (p *BasePkg) Atom() *Cpv {
	ptr := C.pkgcraft_pkg_atom(p.ptr)
	return &Cpv{ptr: ptr}
}

// Return a package's EAPI.
func (p *BasePkg) Eapi() *Eapi {
	ptr := C.pkgcraft_pkg_eapi(p.ptr)
	s := C.pkgcraft_eapi_as_str(ptr)
	defer C.pkgcraft_str_free(s)
	return EAPIS[C.GoString(s)]
}

// Return a package's repo.
func (p *BasePkg) Repo() *BaseRepo {
	ptr := C.pkgcraft_pkg_repo(p.ptr)
	return repoFromPtr(ptr)
}

// Return a package's version.
func (p *BasePkg) Version() *Version {
	ptr := C.pkgcraft_pkg_version(p.ptr)
	return &Version{ptr}
}

func (p *BasePkg) String() string {
	s := C.pkgcraft_pkg_str(p.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
