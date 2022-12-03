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

func (r *EbuildRepo) createPkg(ptr *C.Pkg) *EbuildPkg {
	format := PkgFormat(C.pkgcraft_pkg_format(ptr))
	pkg := &EbuildPkg{&BasePkg{ptr, format}}
	runtime.SetFinalizer(pkg, func(p *EbuildPkg) { C.pkgcraft_pkg_free(p.ptr) })
	return pkg
}

// Return an iterator over the packages of a repo.
func (r *EbuildRepo) PkgIter() *pkgIter[*EbuildPkg] {
	return newPkgIter[*EbuildPkg](r)
}

// Return a channel iterating over the packages of a repo.
func (r *EbuildRepo) Pkgs() <-chan *EbuildPkg {
	return repoPkgs[*EbuildPkg](r)
}

// Return an iterator over the restricted packages of a repo.
func (r *EbuildRepo) RestrictPkgIter(restrict *Restrict) *restrictPkgIter[*EbuildPkg] {
	return newRestrictPkgIter[*EbuildPkg](r, restrict)
}

// Return a channel iterating over the restricted packages of a repo.
func (r *EbuildRepo) RestrictPkgs(restrict *Restrict) <-chan *EbuildPkg {
	return repoRestrictPkgs[*EbuildPkg](r, restrict)
}
