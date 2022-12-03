package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
	"unsafe"
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

// Return an iterator over the packages of a repo.
func (r *BaseRepo) PkgIter() *pkgIter[*BasePkg] {
	return newPkgIter[*BasePkg](r)
}

// Return a channel iterating over the packages of a repo.
func (r *BaseRepo) Pkgs() <-chan *BasePkg {
	return repoPkgs[*BasePkg](r)
}

// Return an iterator over the restricted packages of a repo.
func (r *BaseRepo) RestrictPkgIter(restrict *Restrict) *restrictPkgIter[*BasePkg] {
	return newRestrictPkgIter[*BasePkg](r, restrict)
}

// Return a channel iterating over the restricted packages of a repo.
func (r *BaseRepo) RestrictPkgs(restrict *Restrict) <-chan *BasePkg {
	return repoRestrictPkgs[*BasePkg](r, restrict)
}

// Return true if a repo contains a given object, false otherwise.
func (r *BaseRepo) Contains(obj interface{}) bool {
	switch obj := obj.(type) {
	case string:
		c_str := C.CString(obj)
		defer C.free(unsafe.Pointer(c_str))
		return bool(C.pkgcraft_repo_contains_path(r.ptr, c_str))
	case *Restrict:
		pkgs := r.RestrictPkgs(obj)
		_, ok := <-pkgs
		return ok
	default:
		if restrict, _ := NewRestrict(obj); restrict != nil {
			return r.Contains(restrict)
		}
		return false
	}
}

// Return a new repo from a given pointer.
func repoFromPtr(ptr *C.Repo) *BaseRepo {
	format := RepoFormat(C.pkgcraft_repo_format(ptr))
	return &BaseRepo{ptr, format}
}
