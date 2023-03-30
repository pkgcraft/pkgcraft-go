package pkgcraft_test

import (
	"os"

	"github.com/pelletier/go-toml"
)

type validDep struct {
	Dep      string
	Eapis    string
	Category string
	Package  string
	Blocker  string
	Version  string
	Revision string
	Slot     string
	Subslot  string
	Slot_Op  string
	Use      []string
}

type intersectsDep struct {
	Vals   []string
	Status bool
}

type sortedDep struct {
	Sorted []string
	Equal  bool
}

type depData struct {
	Valid      []validDep
	Invalid    []string
	Intersects []intersectsDep
	Sorting    []sortedDep
}

func parseDepToml() depData {
	var data depData
	f, err := os.ReadFile("testdata/toml/dep.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(f, &data)
	if err != nil {
		panic(err)
	}
	return data
}

var DEP_TOML = parseDepToml()

type intersectsVersion struct {
	Vals   []string
	Status bool
}

type sortedVersion struct {
	Sorted []string
	Equal  bool
}

type hashingVersion struct {
	Versions []string
	Equal    bool
}

type versionData struct {
	Compares   []string
	Intersects []intersectsVersion
	Sorting    []sortedVersion
	Hashing    []hashingVersion
}

func parseVersionToml() versionData {
	var data versionData
	f, err := os.ReadFile("testdata/toml/version.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(f, &data)
	if err != nil {
		panic(err)
	}
	return data
}

var VERSION_TOML = parseVersionToml()
