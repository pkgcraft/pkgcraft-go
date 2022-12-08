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

type pkgIter[P Pkg] struct {
	ptr  *C.RepoPkgIter
	repo pkgRepo[P]
	next P
}

// Create an iterator over the packages of a repo.
func newPkgIter[P Pkg](repo pkgRepo[P]) *pkgIter[P] {
	ptr := C.pkgcraft_repo_iter(repo.p())
	iter := &pkgIter[P]{ptr: ptr, repo: repo}
	runtime.SetFinalizer(iter, func(self *pkgIter[P]) { C.pkgcraft_repo_iter_free(self.ptr) })
	return iter
}

// Determine if a package iterator has another entry.
func (self *pkgIter[P]) HasNext() bool {
	ptr := C.pkgcraft_repo_iter_next(self.ptr)
	if ptr != nil {
		self.next = self.repo.createPkg(ptr)
		return true
	} else {
		return false
	}
}

// Return the next available package in the iterator.
func (self *pkgIter[P]) Next() P {
	return self.next
}

// Return a generic channel iterating over the packages of a repo.
func repoPkgs[P Pkg](repo pkgRepo[P]) <-chan P {
	pkgs := make(chan P)
	go func() {
		for iter := newPkgIter[P](repo); iter.HasNext(); {
			pkgs <- iter.Next()
		}
		close(pkgs)
	}()
	return pkgs
}

type restrictPkgIter[P Pkg] struct {
	ptr  *C.RepoRestrictPkgIter
	repo pkgRepo[P]
	next P
}

// Create a restricted iterator over the packages of a repo.
func newRestrictPkgIter[P Pkg](repo pkgRepo[P], restrict *Restrict) *restrictPkgIter[P] {
	ptr := C.pkgcraft_repo_restrict_iter(repo.p(), restrict.ptr)
	iter := &restrictPkgIter[P]{ptr: ptr, repo: repo}
	runtime.SetFinalizer(iter, func(self *restrictPkgIter[P]) { C.pkgcraft_repo_restrict_iter_free(self.ptr) })
	return iter
}

// Determine if a restricted package iterator has another entry.
func (self *restrictPkgIter[P]) HasNext() bool {
	ptr := C.pkgcraft_repo_restrict_iter_next(self.ptr)
	if ptr != nil {
		self.next = self.repo.createPkg(ptr)
		return true
	} else {
		return false
	}
}

// Return the next available package in the iterator.
func (self *restrictPkgIter[P]) Next() P {
	return self.next
}

// Return a generic channel iterating over the restricted packages of a repo.
func repoRestrictPkgs[P Pkg](repo pkgRepo[P], restrict *Restrict) <-chan P {
	pkgs := make(chan P)

	go func(restrict *Restrict) {
		for iter := newRestrictPkgIter[P](repo, restrict); iter.HasNext(); {
			pkgs <- iter.Next()
		}
		close(pkgs)
	}(restrict)

	return pkgs
}
