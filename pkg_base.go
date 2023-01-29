package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type BasePkg struct {
	ptr    *C.Pkg
	format PkgFormat
	// cached fields
	eapi *Eapi
}

func (self *BasePkg) p() *C.Pkg {
	return self.ptr
}

// Return a package's atom.
func (self *BasePkg) Cpv() *Cpv {
	ptr := C.pkgcraft_pkg_cpv(self.ptr)
	cpv, _ := cpvFromPtr(ptr)
	return cpv
}

// Return a package's EAPI.
func (self *BasePkg) Eapi() *Eapi {
	if self.eapi == nil {
		ptr := C.pkgcraft_pkg_eapi(self.ptr)
		s := C.pkgcraft_eapi_as_str(ptr)
		defer C.pkgcraft_str_free(s)
		self.eapi = EAPIS[C.GoString(s)]
	}
	return self.eapi
}

// Return a package's repo.
func (self *BasePkg) Repo() *BaseRepo {
	ptr := C.pkgcraft_pkg_repo(self.ptr)
	return repoFromPtr(ptr)
}

// Return a package's version.
func (self *BasePkg) Version() *Version {
	ptr := C.pkgcraft_pkg_version(self.ptr)
	version, _ := versionFromPtr(ptr)
	return version
}

func (self *BasePkg) String() string {
	s := C.pkgcraft_pkg_str(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Compare a package with another package returning -1, 0, or 1 if the first is
// less than, equal to, or greater than the second, respectively.
func (self *BasePkg) Cmp(other pkgPtr) int {
	return int(C.pkgcraft_pkg_cmp(self.ptr, other.p()))
}
