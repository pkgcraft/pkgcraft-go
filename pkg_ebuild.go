package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type EbuildPkg struct {
	*BasePkg
}

// Return a package's repo.
func (self *EbuildPkg) Repo() *EbuildRepo {
	base := &BaseRepo{C.pkgcraft_pkg_repo(self.ptr), RepoFormatEbuild}
	return &EbuildRepo{base, nil}
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
		return "", newPkgcraftError()
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

// Return a package's dependencies for the given descriptors.
func (self *EbuildPkg) Dependencies(keys []string) (*DependencySet, error) {
	c_keys, c_len := sliceToCharArray(keys)
	ptr := C.pkgcraft_pkg_ebuild_dependencies(self.ptr, c_keys, c_len)
	if ptr != nil {
		return dependencySetFromPtr(ptr), nil
	} else {
		return nil, newPkgcraftError()
	}
}

// Return a package's DEPEND.
func (self *EbuildPkg) Depend() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_depend(self.ptr))
}

// Return a package's BDEPEND.
func (self *EbuildPkg) Bdepend() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_bdepend(self.ptr))
}

// Return a package's IDEPEND.
func (self *EbuildPkg) Idepend() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_idepend(self.ptr))
}

// Return a package's PDEPEND.
func (self *EbuildPkg) Pdepend() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_pdepend(self.ptr))
}

// Return a package's RDEPEND.
func (self *EbuildPkg) Rdepend() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_rdepend(self.ptr))
}

// Return a package's LICENSE.
func (self *EbuildPkg) License() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_license(self.ptr))
}

// Return a package's PROPERTIES.
func (self *EbuildPkg) Properties() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_properties(self.ptr))
}

// Return a package's REQUIRED_USE.
func (self *EbuildPkg) RequiredUse() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_required_use(self.ptr))
}

// Return a package's Restrict.
func (self *EbuildPkg) Restrict() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_restrict(self.ptr))
}

// Return a package's SRC_URI.
func (self *EbuildPkg) SrcUri() *DependencySet {
	return dependencySetFromPtr(C.pkgcraft_pkg_ebuild_src_uri(self.ptr))
}
