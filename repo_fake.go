package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type FakeRepo struct {
	*BaseRepo
}

// Create a new fake repo.
func NewFakeRepo(id string, priority int, cpvs []string) (*FakeRepo, error) {
	c_cpvs := C.malloc(C.size_t(len(cpvs)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	a := (*[1<<30 - 1]*C.char)(c_cpvs)
	for i, s := range cpvs {
		a[i] = C.CString(s)
	}

	c_id := C.CString(id)
	ptr := C.pkgcraft_repo_fake_new(c_id, C.int(priority), (**C.char)(c_cpvs), C.size_t(len(cpvs)))
	C.free(unsafe.Pointer(c_id))
	C.free(c_cpvs)

	if ptr != nil {
		repo := &FakeRepo{repoFromPtr(ptr)}
		runtime.SetFinalizer(repo, func(r *FakeRepo) { C.pkgcraft_repo_free(r.ptr) })
		return repo, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

func (r *FakeRepo) createPkg(ptr *C.Pkg) *FakePkg {
	format := PkgFormat(C.pkgcraft_pkg_format(ptr))
	pkg := &FakePkg{&BasePkg{ptr, format}}
	runtime.SetFinalizer(pkg, func(p *FakePkg) { C.pkgcraft_pkg_free(p.ptr) })
	return pkg
}

// Return a channel iterating over the packages of a repo.
func (r *FakeRepo) Pkgs() <-chan *FakePkg {
	return repoPkgs((pkgRepo[*FakePkg])(r))
}

// Return a channel iterating over the restricted packages of a repo.
func (r *FakeRepo) RestrictPkgs(restrict *Restrict) <-chan *FakePkg {
	return repoRestrictPkgs((pkgRepo[*FakePkg])(r), restrict)
}
