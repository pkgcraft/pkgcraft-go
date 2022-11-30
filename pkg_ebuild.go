package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
)

type EbuildPkg struct {
	*BasePkg
}

// Return a package's repo.
func (p *EbuildPkg) Repo() *EbuildRepo {
	base := &BaseRepo{C.pkgcraft_pkg_repo(p.ptr), RepoFormatEbuild}
	return &EbuildRepo{base}
}

// Return a package's path.
func (p *EbuildPkg) Path() string {
	s := C.pkgcraft_pkg_ebuild_path(p.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package's ebuild file content.
func (p *EbuildPkg) Ebuild() (string, error) {
	s := C.pkgcraft_pkg_ebuild_ebuild(p.ptr)
	if s != nil {
		defer C.pkgcraft_str_free(s)
		return C.GoString(s), nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return "", errors.New(C.GoString(s))
	}
}

// Return a package's description.
func (p *EbuildPkg) Description() string {
	s := C.pkgcraft_pkg_ebuild_description(p.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package's slot.
func (p *EbuildPkg) Slot() string {
	s := C.pkgcraft_pkg_ebuild_slot(p.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package's subslot.
func (p *EbuildPkg) Subslot() string {
	s := C.pkgcraft_pkg_ebuild_subslot(p.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}
