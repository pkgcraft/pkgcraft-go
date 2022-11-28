package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type Repo interface {
	id() string
	path() string
	is_empty() bool
}

type BaseRepo struct {
	ptr *C.Repo
}

// Return a repo's id.
func (r *BaseRepo) id() string {
	s := C.pkgcraft_repo_id(r.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a repo's path.
func (r *BaseRepo) path() string {
	s := C.pkgcraft_repo_path(r.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return if a repo is empty.
func (r *BaseRepo) is_empty() bool {
	return bool(C.pkgcraft_repo_is_empty(r.ptr))
}

type EbuildRepo struct {
	*BaseRepo
}

type FakeRepo struct {
	*BaseRepo
}

type RepoFormat int

const (
	RepoFormatEbuild RepoFormat = iota
	RepoFormatFake
	RepoFormatEmpty
)

// Return a new repo from a given pointer.
func repo_from_ptr(r *C.Repo) (Repo) {
	var repo Repo

	base := &BaseRepo{ptr: r}
	runtime.SetFinalizer(base, func(r *BaseRepo) { C.pkgcraft_repo_free(r.ptr) })

	format := RepoFormat(C.pkgcraft_repo_format(r))
	if format == RepoFormatEbuild {
		repo = &EbuildRepo{base}
	} else if format == RepoFormatFake {
		repo = &FakeRepo{base}
	} else {
		panic("unsupported repo format")
	}

	return repo
}
