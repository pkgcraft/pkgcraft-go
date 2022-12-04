package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"fmt"
	"unsafe"

	. "github.com/pkgcraft/pkgcraft-go"
)

func BlockerFromString(s string) (Blocker, error) {
	c_str := C.CString(s)
	i := C.pkgcraft_atom_blocker_from_str(c_str)
	C.free(unsafe.Pointer(c_str))
	if i >= 0 {
		return Blocker(i), nil
	} else {
		return BlockerNone, fmt.Errorf("invalid blocker: %s", s)
	}
}

func SlotOperatorFromString(s string) (SlotOperator, error) {
	c_str := C.CString(s)
	i := C.pkgcraft_atom_slot_op_from_str(c_str)
	C.free(unsafe.Pointer(c_str))
	if i >= 0 {
		return SlotOperator(i), nil
	} else {
		return SlotOpNone, fmt.Errorf("invalid slot operator: %s", s)
	}
}
