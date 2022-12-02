package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type pkgRepo[P Pkg] interface {
	p() *C.Repo
	createPkg(*C.Pkg) P
}

// Return a generic channel iterating over the packages of a repo.
func repoPkgs[P Pkg](r pkgRepo[P]) <-chan P {
	pkgs := make(chan P)

	go func() {
		iter := C.pkgcraft_repo_iter(r.p())
		for {
			ptr := C.pkgcraft_repo_iter_next(iter)
			if ptr != nil {
				pkgs <- r.createPkg(ptr)
			} else {
				break
			}
		}
		close(pkgs)
		C.pkgcraft_repo_iter_free(iter)
	}()

	return pkgs
}

type RepoFormat int

const (
	RepoFormatEbuild RepoFormat = iota
	RepoFormatFake
)
