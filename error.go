package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

type PkgcraftError struct {
	msg string
}

func (self *PkgcraftError) Error() string {
	return self.msg
}

func newPkgcraftError() error {
	err := C.pkgcraft_error_last()
	if err != nil {
		defer C.pkgcraft_error_free(err)
		return &PkgcraftError{C.GoString(err.message)}
	} else {
		panic("no pkgcraft error occurred")
	}
}
