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

// Return a package's DEPEND.
func (p *EbuildPkg) Depend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_depend(p.ptr), DepSetAtom)
}

// Return a package's BDEPEND.
func (p *EbuildPkg) Bdepend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_bdepend(p.ptr), DepSetAtom)
}

// Return a package's IDEPEND.
func (p *EbuildPkg) Idepend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_idepend(p.ptr), DepSetAtom)
}

// Return a package's PDEPEND.
func (p *EbuildPkg) Pdepend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_pdepend(p.ptr), DepSetAtom)
}

// Return a package's RDEPEND.
func (p *EbuildPkg) Rdepend() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_rdepend(p.ptr), DepSetAtom)
}

// Return a package's LICENSE.
func (p *EbuildPkg) License() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_license(p.ptr), DepSetString)
}

// Return a package's PROPERTIES.
func (p *EbuildPkg) Properties() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_properties(p.ptr), DepSetString)
}

// Return a package's REQUIRED_USE.
func (p *EbuildPkg) RequiredUse() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_required_use(p.ptr), DepSetString)
}

// Return a package's Restrict.
func (p *EbuildPkg) Restrict() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_restrict(p.ptr), DepSetString)
}

// Return a package's SRC_URI.
func (p *EbuildPkg) SrcUri() *DepSet {
	return depSetFromPtr(C.pkgcraft_pkg_ebuild_src_uri(p.ptr), DepSetUri)
}
