package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type BasePkg struct {
	ptr    *C.Pkg
	format PkgFormat
	// cached fields
	eapi *Eapi
	cpv  *Cpv
}

func (self *BasePkg) p() *C.Pkg {
	return self.ptr
}

// Return a package's package and version.
func (self *BasePkg) P() string {
	return self.Cpv().P()
}

// Return a package's package, version, and revision.
func (self *BasePkg) Pf() string {
	return self.Cpv().Pf()
}

// Return an package's revision.
func (self *BasePkg) Pr() string {
	return self.Cpv().Pr()
}

// Return an package's version.
func (self *BasePkg) Pv() string {
	return self.Cpv().Pv()
}

// Return a package's version and revision.
func (self *BasePkg) Pvr() string {
	return self.Cpv().Pvr()
}

// Return a package's Cpv.
func (self *BasePkg) Cpv() *Cpv {
	if self.cpv == nil {
		ptr := C.pkgcraft_pkg_cpv(self.ptr)
		self.cpv, _ = cpvFromPtr(ptr)
	}
	return self.cpv
}

// Return a package's Cpn.
func (self *BasePkg) Cpn() *Cpn {
	return self.Cpv().Cpn()
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
	repo := repoFromPtr(ptr)
	runtime.SetFinalizer(repo, func(self *BaseRepo) { C.pkgcraft_repo_free(self.ptr) })
	return repo
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
