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

func (self *BaseRepo) p() *C.Repo {
	return self.ptr
}

// Return a repo's id.
func (self *BaseRepo) Id() string {
	s := C.pkgcraft_repo_id(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a repo's path.
func (self *BaseRepo) Path() string {
	s := C.pkgcraft_repo_path(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return if a repo is empty.
func (self *BaseRepo) IsEmpty() bool {
	return bool(C.pkgcraft_repo_is_empty(self.ptr))
}

// Return the number of packages in a repo.
func (self *BaseRepo) Len() int {
	return int(C.pkgcraft_repo_len(self.ptr))
}

func (self *BaseRepo) String() string {
	return self.Id()
}

// Compare a repo with another repo returning -1, 0, or 1 if the first is less
// than, equal to, or greater than the second, respectively.
func (self *BaseRepo) Cmp(other repoPtr) int {
	return int(C.pkgcraft_repo_cmp(self.ptr, other.p()))
}

func (self *BaseRepo) createPkg(ptr *C.Pkg) *BasePkg {
	format := PkgFormat(C.pkgcraft_pkg_format(ptr))
	pkg := &BasePkg{ptr, format}
	runtime.SetFinalizer(pkg, func(self *BasePkg) { C.pkgcraft_pkg_free(self.ptr) })
	return pkg
}

// Return an iterator over the packages of a repo.
func (self *BaseRepo) PkgIter() *pkgIter[*BasePkg] {
	return newPkgIter[*BasePkg](self)
}

// Return a channel iterating over the packages of a repo.
func (self *BaseRepo) Pkgs() <-chan *BasePkg {
	return repoPkgs[*BasePkg](self)
}

// Return an iterator over the restricted packages of a repo.
func (self *BaseRepo) RestrictPkgIter(restrict *Restrict) *restrictPkgIter[*BasePkg] {
	return newRestrictPkgIter[*BasePkg](self, restrict)
}

// Return a channel iterating over the restricted packages of a repo.
func (self *BaseRepo) RestrictPkgs(restrict *Restrict) <-chan *BasePkg {
	return repoRestrictPkgs[*BasePkg](self, restrict)
}

// Return true if a repo contains a given object, false otherwise.
func (self *BaseRepo) Contains(obj interface{}) bool {
	switch obj := obj.(type) {
	case string:
		c_str := C.CString(obj)
		defer C.free(unsafe.Pointer(c_str))
		return bool(C.pkgcraft_repo_contains_path(self.ptr, c_str))
	case *Restrict:
		pkgs := self.RestrictPkgs(obj)
		_, ok := <-pkgs
		return ok
	default:
		if restrict, _ := NewRestrict(obj); restrict != nil {
			return self.Contains(restrict)
		}
		return false
	}
}

// Return a new repo from a given pointer.
func repoFromPtr(ptr *C.Repo) *BaseRepo {
	format := RepoFormat(C.pkgcraft_repo_format(ptr))
	return &BaseRepo{ptr, format}
}
