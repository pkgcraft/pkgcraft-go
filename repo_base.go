package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type BaseRepo struct {
	ptr    *C.Repo
	format RepoFormat
}

func (r *BaseRepo) p() *C.Repo {
	return r.ptr
}

// Return a repo's id.
func (r *BaseRepo) Id() string {
	s := C.pkgcraft_repo_id(r.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a repo's path.
func (r *BaseRepo) Path() string {
	s := C.pkgcraft_repo_path(r.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return if a repo is empty.
func (r *BaseRepo) IsEmpty() bool {
	return bool(C.pkgcraft_repo_is_empty(r.ptr))
}

// Return the number of packages in a repo.
func (r *BaseRepo) Len() int {
	return int(C.pkgcraft_repo_len(r.ptr))
}

func (r *BaseRepo) String() string {
	return r.Id()
}

// Compare a repo with another repo returning -1, 0, or 1 if the first is less
// than, equal to, or greater than the second, respectively.
func (r1 *BaseRepo) Cmp(r2 *BaseRepo) int {
	return int(C.pkgcraft_repo_cmp(r1.ptr, r2.ptr))
}

func (r *BaseRepo) createPkg(ptr *C.Pkg) *BasePkg {
	format := PkgFormat(C.pkgcraft_pkg_format(ptr))
	pkg := &BasePkg{ptr, format}
	runtime.SetFinalizer(pkg, func(p *BasePkg) { C.pkgcraft_pkg_free(p.ptr) })
	return pkg
}

// Return a channel iterating over the packages of a repo.
func (r *BaseRepo) Pkgs() <-chan *BasePkg {
	return repoPkgs((pkgRepo[*BasePkg])(r))
}

// Return a new repo from a given pointer.
func repoFromPtr(ptr *C.Repo) *BaseRepo {
	format := RepoFormat(C.pkgcraft_repo_format(ptr))
	return &BaseRepo{ptr, format}
}
