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
	defer C.free(unsafe.Pointer(c_id))
	ptr := C.pkgcraft_repo_fake_new(c_id, C.int(priority), (**C.char)(c_cpvs), C.size_t(len(cpvs)))
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

// Return a channel iterating over the packages of a repo.
func (r *FakeRepo) Pkgs() <-chan *FakePkg {
	pkgs := make(chan *FakePkg)

	go func() {
		iter := C.pkgcraft_repo_iter(r.ptr)
		for {
			ptr := C.pkgcraft_repo_iter_next(iter)
			if ptr != nil {
				pkgs <- &FakePkg{pkgFromPtr(ptr)}
			} else {
				break
			}
		}
		close(pkgs)
		C.pkgcraft_repo_iter_free(iter)
	}()

	return pkgs
}
