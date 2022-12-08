package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type Pkg interface {
	Atom() *Cpv
	Eapi() *Eapi
	Version() *Version
	String() string
}

type PkgFormat int

const (
	PkgFormatEbuild PkgFormat = iota
	PkgFormatFake
)

type pkgPtr interface {
	p() *C.Pkg
}
