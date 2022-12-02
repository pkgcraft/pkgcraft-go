package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type pkgRepo[P Pkg] interface {
	p() *C.Repo
	createPkg(*C.Pkg) P
}

// Return a generic channel iterating over the packages of a repo.
func repoPkgs[P Pkg](repo pkgRepo[P]) <-chan P {
	pkgs := make(chan P)

	go func() {
		iter := C.pkgcraft_repo_iter(repo.p())
		for {
			ptr := C.pkgcraft_repo_iter_next(iter)
			if ptr != nil {
				pkgs <- repo.createPkg(ptr)
			} else {
				break
			}
		}
		close(pkgs)
		C.pkgcraft_repo_iter_free(iter)
	}()

	return pkgs
}

// Return a generic channel iterating over the restricted packages of a repo.
func repoRestrictPkgs[P Pkg](repo pkgRepo[P], restrict *Restrict) <-chan P {
	pkgs := make(chan P)

	go func(restrict *Restrict) {
		iter := C.pkgcraft_repo_restrict_iter(repo.p(), restrict.ptr)
		for {
			ptr := C.pkgcraft_repo_restrict_iter_next(iter)
			if ptr != nil {
				pkgs <- repo.createPkg(ptr)
			} else {
				break
			}
		}
		close(pkgs)
		C.pkgcraft_repo_restrict_iter_free(iter)
	}(restrict)

	return pkgs
}

type RepoFormat int

const (
	RepoFormatEbuild RepoFormat = iota
	RepoFormatFake
)
