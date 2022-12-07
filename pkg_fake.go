package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type FakePkg struct {
	*BasePkg
}

// Return a package's repo.
func (self *FakePkg) Repo() *FakeRepo {
	base := &BaseRepo{C.pkgcraft_pkg_repo(self.ptr), RepoFormatFake}
	return &FakeRepo{base}
}
