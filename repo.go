package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"runtime"
)

type RepoFormat int

const (
	RepoFormatEbuild RepoFormat = iota
	RepoFormatFake
)

type repoPtr interface {
	p() *C.Repo
}

type pkgRepo[P Pkg] interface {
	repoPtr
	createPkg(*C.Pkg) P
}

type repoIter[P Pkg] struct {
	ptr  *C.RepoIter
	repo pkgRepo[P]
	next P
}

// Create an iterator over the packages of a repo.
func newRepoIter[P Pkg](repo pkgRepo[P]) *repoIter[P] {
	ptr := C.pkgcraft_repo_iter(repo.p())
	iter := &repoIter[P]{ptr: ptr, repo: repo}
	runtime.SetFinalizer(iter, func(self *repoIter[P]) { C.pkgcraft_repo_iter_free(self.ptr) })
	return iter
}

// Determine if a package iterator has another entry.
func (self *repoIter[P]) HasNext() bool {
	ptr := C.pkgcraft_repo_iter_next(self.ptr)
	if ptr != nil {
		self.next = self.repo.createPkg(ptr)
		return true
	} else {
		return false
	}
}

// Return the next available package in the iterator.
func (self *repoIter[P]) Next() P {
	return self.next
}

// Return a generic channel iterating over the packages of a repo.
func repoPkgs[P Pkg](repo pkgRepo[P]) <-chan P {
	pkgs := make(chan P)
	go func() {
		for iter := newRepoIter[P](repo); iter.HasNext(); {
			pkgs <- iter.Next()
		}
		close(pkgs)
	}()
	return pkgs
}

type repoIterRestrict[P Pkg] struct {
	ptr  *C.RepoIterRestrict
	repo pkgRepo[P]
	next P
}

// Create a restricted iterator over the packages of a repo.
func newRepoIterRestrict[P Pkg](repo pkgRepo[P], restrict *Restrict) *repoIterRestrict[P] {
	ptr := C.pkgcraft_repo_iter_restrict(repo.p(), restrict.ptr)
	iter := &repoIterRestrict[P]{ptr: ptr, repo: repo}
	runtime.SetFinalizer(iter, func(self *repoIterRestrict[P]) { C.pkgcraft_repo_iter_restrict_free(self.ptr) })
	return iter
}

// Determine if a restricted package iterator has another entry.
func (self *repoIterRestrict[P]) HasNext() bool {
	ptr := C.pkgcraft_repo_iter_restrict_next(self.ptr)
	if ptr != nil {
		self.next = self.repo.createPkg(ptr)
		return true
	} else {
		return false
	}
}

// Return the next available package in the iterator.
func (self *repoIterRestrict[P]) Next() P {
	return self.next
}

// Return a generic channel iterating over the restricted packages of a repo.
func repoRestrictPkgs[P Pkg](repo pkgRepo[P], restrict *Restrict) <-chan P {
	pkgs := make(chan P)

	go func(restrict *Restrict) {
		for iter := newRepoIterRestrict[P](repo, restrict); iter.HasNext(); {
			pkgs <- iter.Next()
		}
		close(pkgs)
	}(restrict)

	return pkgs
}
