package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type BasePkg struct {
	ptr    *C.Pkg
	format PkgFormat
}

// Return a package's atom.
func (self *BasePkg) Atom() *Cpv {
	ptr := C.pkgcraft_pkg_atom(self.ptr)
	return &Cpv{ptr: ptr}
}

// Return a package's EAPI.
func (self *BasePkg) Eapi() *Eapi {
	ptr := C.pkgcraft_pkg_eapi(self.ptr)
	s := C.pkgcraft_eapi_as_str(ptr)
	defer C.pkgcraft_str_free(s)
	return EAPIS[C.GoString(s)]
}

// Return a package's repo.
func (self *BasePkg) Repo() *BaseRepo {
	ptr := C.pkgcraft_pkg_repo(self.ptr)
	return repoFromPtr(ptr)
}

// Return a package's version.
func (self *BasePkg) Version() *Version {
	ptr := C.pkgcraft_pkg_version(self.ptr)
	return &Version{ptr}
}

func (self *BasePkg) String() string {
	s := C.pkgcraft_pkg_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
