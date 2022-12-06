package pkgcraft

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
