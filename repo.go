package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type Repo interface {
	Id() string
	Path() string
	IsEmpty() bool
	String() string
}

type BaseRepo struct {
	ptr *C.Repo
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

type repo_pkg_iter struct {
	ptr *C.RepoPkgIter
	next *BasePkg
}

// Return true if the repo iterator contains another package, false otherwise.
func (i *repo_pkg_iter) HasNext() bool {
	ptr := C.pkgcraft_repo_iter_next(i.ptr)
	if ptr == nil {
		i.next = nil
		return false
	} else {
		i.next = pkg_from_ptr(ptr)
		return true
	}
}

// Return the next package in the iterator.
func (i *repo_pkg_iter) Next() *BasePkg {
	return i.next
}

// Return a new package iterator for a repo.
func (r *BaseRepo) NewPkgIterator() Iterator[*BasePkg] {
	ptr := C.pkgcraft_repo_iter(r.ptr)
	iter := &repo_pkg_iter{ptr, nil}
	runtime.SetFinalizer(iter, func(i *repo_pkg_iter) { C.pkgcraft_repo_iter_free(i.ptr) })
	return iter
}

type FakeRepo struct {
	*BaseRepo
}

type RepoFormat int

const (
	RepoFormatEbuild RepoFormat = iota
	RepoFormatFake
)

// Return a new repo from a given pointer.
func repo_from_ptr(ptr *C.Repo, ref bool) *BaseRepo {
	format := RepoFormat(C.pkgcraft_repo_format(ptr))
	base := &BaseRepo{ptr, format}
	if !ref {
		runtime.SetFinalizer(base, func(r *BaseRepo) { C.pkgcraft_repo_free(r.ptr) })
	}
	return base
}
