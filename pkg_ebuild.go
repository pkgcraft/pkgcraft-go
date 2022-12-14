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
func (self *EbuildPkg) Repo() *EbuildRepo {
	base := &BaseRepo{C.pkgcraft_pkg_repo(self.ptr), RepoFormatEbuild}
	return &EbuildRepo{base}
}

// Return a package's path.
func (self *EbuildPkg) Path() string {
	s := C.pkgcraft_pkg_ebuild_path(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package's ebuild file content.
func (self *EbuildPkg) Ebuild() (string, error) {
	s := C.pkgcraft_pkg_ebuild_ebuild(self.ptr)
	if s != nil {
		defer C.pkgcraft_str_free(s)
		return C.GoString(s), nil
	} else {
		err := C.pkgcraft_error_last()
		defer C.pkgcraft_error_free(err)
		return "", errors.New(C.GoString(err.message))
	}
}

// Return a package's description.
func (self *EbuildPkg) Description() string {
	s := C.pkgcraft_pkg_ebuild_description(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package's slot.
func (self *EbuildPkg) Slot() string {
	s := C.pkgcraft_pkg_ebuild_slot(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package's subslot.
func (self *EbuildPkg) Subslot() string {
	s := C.pkgcraft_pkg_ebuild_subslot(self.ptr)
	defer C.pkgcraft_str_free(s)
	return C.GoString(s)
}

// Return a package's DEPEND.
func (self *EbuildPkg) Depend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_depend(self.ptr), DepSetAtom)
}

// Return a package's BDEPEND.
func (self *EbuildPkg) Bdepend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_bdepend(self.ptr), DepSetAtom)
}

// Return a package's IDEPEND.
func (self *EbuildPkg) Idepend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_idepend(self.ptr), DepSetAtom)
}

// Return a package's PDEPEND.
func (self *EbuildPkg) Pdepend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_pdepend(self.ptr), DepSetAtom)
}

// Return a package's RDEPEND.
func (self *EbuildPkg) Rdepend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_rdepend(self.ptr), DepSetAtom)
}

// Return a package's LICENSE.
func (self *EbuildPkg) License() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_license(self.ptr), DepSetString)
}

// Return a package's PROPERTIES.
func (self *EbuildPkg) Properties() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_properties(self.ptr), DepSetString)
}

// Return a package's REQUIRED_USE.
func (self *EbuildPkg) RequiredUse() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_required_use(self.ptr), DepSetString)
}

// Return a package's Restrict.
func (self *EbuildPkg) Restrict() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_restrict(self.ptr), DepSetString)
}

// Return a package's SRC_URI.
func (self *EbuildPkg) SrcUri() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_src_uri(self.ptr), DepSetUri)
}
