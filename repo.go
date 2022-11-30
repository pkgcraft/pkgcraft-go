package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type Repo interface {
	Id() string
	Path() string
	IsEmpty() bool
	String() string
	Pkgs() <-chan *BasePkg
}

type BaseRepo struct {
	ptr    *C.Repo
	format RepoFormat
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

func (r *BaseRepo) String() string {
	return r.Id()
}

// Return a channel iterating over the packages of a repo.
func (r *BaseRepo) Pkgs() <-chan *BasePkg {
	pkgs := make(chan *BasePkg)

	go func() {
		iter := C.pkgcraft_repo_iter(r.ptr)
		for {
			ptr := C.pkgcraft_repo_iter_next(iter)
			if ptr != nil {
				pkgs <- pkg_from_ptr(ptr)
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

// Return a new repo from a given pointer.
func repo_from_ptr(ptr *C.Repo) *BaseRepo {
	format := RepoFormat(C.pkgcraft_repo_format(ptr))
	return &BaseRepo{ptr, format}
}
