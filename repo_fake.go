package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
	"unsafe"
)

type FakeRepo struct {
	*BaseRepo
}

// Create a new fake repo.
func NewFakeRepo(id string, priority int, cpvs []string) (*FakeRepo, error) {
	c_cpvs, c_len := sliceToCharArray(cpvs)
	c_id := C.CString(id)
	ptr := C.pkgcraft_repo_fake_new(c_id, C.int(priority), c_cpvs, c_len)
	C.free(unsafe.Pointer(c_id))
	C.free(unsafe.Pointer(c_cpvs))

	if ptr != nil {
		repo := &FakeRepo{repoFromPtr(ptr)}
		runtime.SetFinalizer(repo, func(self *FakeRepo) { C.pkgcraft_repo_free(self.ptr) })
		return repo, nil
	} else {
		return nil, newPkgcraftError()
	}
}

// Add packages to an existing repo.
func (self *FakeRepo) Extend(cpvs []string) error {
	c_cpvs, c_len := sliceToCharArray(cpvs)
	ptr := C.pkgcraft_repo_fake_extend(self.ptr, (**C.char)(c_cpvs), c_len)
	C.free(unsafe.Pointer(c_cpvs))

	if ptr != nil {
		return nil
	} else {
		return newPkgcraftError()
	}
}

func (self *FakeRepo) createPkg(ptr *C.Pkg) *FakePkg {
	format := PkgFormat(C.pkgcraft_pkg_format(ptr))
	pkg := &FakePkg{&BasePkg{ptr: ptr, format: format}}
	runtime.SetFinalizer(pkg, func(self *FakePkg) { C.pkgcraft_pkg_free(self.ptr) })
	return pkg
}

// Return an iterator over the packages of a repo.
func (self *FakeRepo) Iter() *repoIter[*FakePkg] {
	return newRepoIter[*FakePkg](self)
}

// Return a channel iterating over the packages of a repo.
func (self *FakeRepo) Pkgs() <-chan *FakePkg {
	return repoPkgs((pkgRepo[*FakePkg])(self))
}

// Return an iterator over the restricted packages of a repo.
func (self *FakeRepo) IterRestrict(restrict *Restrict) *repoIterRestrict[*FakePkg] {
	return newRepoIterRestrict[*FakePkg](self, restrict)
}

// Return a channel iterating over the restricted packages of a repo.
func (self *FakeRepo) RestrictPkgs(restrict *Restrict) <-chan *FakePkg {
	return repoRestrictPkgs((pkgRepo[*FakePkg])(self), restrict)
}
