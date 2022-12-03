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

// Convert a slice of Go strings to an array of C strings.
func sliceToCharArray(vals []string) (**C.char, C.size_t) {
	c_strs := C.malloc(C.size_t(len(vals)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	a := (*[1<<30 - 1]*C.char)(c_strs)
	for i, s := range vals {
		a[i] = C.CString(s)
	}
	return (**C.char)(c_strs), C.size_t(len(vals))
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
		runtime.SetFinalizer(repo, func(r *FakeRepo) { C.pkgcraft_repo_free(r.ptr) })
		return repo, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

// Add packages to an existing repo.
func (r *FakeRepo) Extend(cpvs []string) error {
	c_cpvs, c_len := sliceToCharArray(cpvs)
	ptr := C.pkgcraft_repo_fake_extend(r.ptr, (**C.char)(c_cpvs), c_len)
	C.free(unsafe.Pointer(c_cpvs))

	if ptr != nil {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return errors.New(C.GoString(s))
	}
	return nil
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
