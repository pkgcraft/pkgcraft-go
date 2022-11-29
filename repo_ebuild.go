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

type repo_ebuild_pkg_iter struct {
	ptr *C.RepoPkgIter
	next *EbuildPkg
}

// Return true if the repo iterator contains another package, false otherwise.
func (i *repo_ebuild_pkg_iter) HasNext() bool {
	ptr := C.pkgcraft_repo_iter_next(i.ptr)
	if ptr == nil {
		i.next = nil
		return false
	} else {
		pkg := pkg_from_ptr(ptr)
		if pkg.format == PkgFormatEbuild {
			i.next = &EbuildPkg{pkg}
		} else {
			panic("invalid pkg format")
		}
		return true
	}
}

// Return the next package in the iterator.
func (i *repo_ebuild_pkg_iter) Next() *EbuildPkg {
	return i.next
}

// Return a new package iterator for a repo.
func (r *EbuildRepo) NewPkgIterator() Iterator[*EbuildPkg] {
	ptr := C.pkgcraft_repo_iter(r.ptr)
	iter := &repo_ebuild_pkg_iter{ptr, nil}
	runtime.SetFinalizer(iter, func(i *repo_ebuild_pkg_iter) { C.pkgcraft_repo_iter_free(i.ptr) })
	return iter
}
