package pkgcraft

// #cgo pkg-config: pkgcraft

import "testing"

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%q != %q", a, b)
	}
}

func TestAtom(t *testing.T) {
	var atom *Atom

	// unversioned
	atom, _ = NewAtom("cat/pkg")
	assertEqual(t, atom.category(), "cat")
	assertEqual(t, atom.pn(), "pkg")
	assertEqual(t, atom.version(), "")

	// versioned
	atom, _ = NewAtom("=cat/pkg-2")
	assertEqual(t, atom.category(), "cat")
	assertEqual(t, atom.pn(), "pkg")
	assertEqual(t, atom.version(), "2")

	// slotted
	atom, _ = NewAtom("cat/pkg:1")
	assertEqual(t, atom.slot(), "1")
	assertEqual(t, atom.subslot(), "")

	// subslotted
	atom, _ = NewAtom("cat/pkg:1/2")
	assertEqual(t, atom.slot(), "1")
	assertEqual(t, atom.subslot(), "2")
}

func BenchmarkNewAtom(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewAtom("=cat/pkg-1-r2:3/4=[a,b,c]")
	}
}

func TestVersion(t *testing.T) {
	var version *Version

	// non-revision
	version, _ = NewVersion("1")
	assertEqual(t, version.revision(), "")

	// revisioned
	version, _ = NewVersion("1-r1")
	assertEqual(t, version.revision(), "1")

	// explicit '0' revision
	version, _ = NewVersion("1-r0")
	assertEqual(t, version.revision(), "0")
}

func BenchmarkNewVersion(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewVersion("1.2.3_alpha4-r5")
	}
}
