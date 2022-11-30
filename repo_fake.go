package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type FakeRepo struct {
	*BaseRepo
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
