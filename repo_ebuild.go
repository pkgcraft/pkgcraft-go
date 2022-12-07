package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type EbuildRepo struct {
	*BaseRepo
}

func (self *EbuildRepo) createPkg(ptr *C.Pkg) *EbuildPkg {
	format := PkgFormat(C.pkgcraft_pkg_format(ptr))
	pkg := &EbuildPkg{&BasePkg{ptr, format}}
	runtime.SetFinalizer(pkg, func(self *EbuildPkg) { C.pkgcraft_pkg_free(self.ptr) })
	return pkg
}

// Return an iterator over the packages of a repo.
func (self *EbuildRepo) PkgIter() *pkgIter[*EbuildPkg] {
	return newPkgIter[*EbuildPkg](self)
}

// Return a channel iterating over the packages of a repo.
func (self *EbuildRepo) Pkgs() <-chan *EbuildPkg {
	return repoPkgs[*EbuildPkg](self)
}

// Return an iterator over the restricted packages of a repo.
func (self *EbuildRepo) RestrictPkgIter(restrict *Restrict) *restrictPkgIter[*EbuildPkg] {
	return newRestrictPkgIter[*EbuildPkg](self, restrict)
}

// Return a channel iterating over the restricted packages of a repo.
func (self *EbuildRepo) RestrictPkgs(restrict *Restrict) <-chan *EbuildPkg {
	return repoRestrictPkgs[*EbuildPkg](self, restrict)
}
