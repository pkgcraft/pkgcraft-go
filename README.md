[![CI](https://github.com/pkgcraft/pkgcraft-go/workflows/CI/badge.svg)](https://github.com/pkgcraft/pkgcraft-go/actions/workflows/ci.yml)
[![coverage](https://codecov.io/gh/pkgcraft/pkgcraft-go/branch/main/graph/badge.svg)](https://codecov.io/gh/pkgcraft/pkgcraft-go)

# pkgcraft-go

Go bindings for pkgcraft.

## Development

Requirements: >=go-1.18 and everything required to build
[pkgcraft-c](https://github.com/pkgcraft/pkgcraft-c)

Use the following commands to set up a dev environment:

```bash
# clone the pkgcraft workspace and pull the latest project updates
git clone --recurse-submodules https://github.com/pkgcraft/pkgcraft-workspace.git
cd pkgcraft-workspace
git submodule update --recursive --remote

# build pkgcraft-c library and set shell variables (e.g. $PKG_CONFIG_PATH)
source ./build pkgcraft-c

cd pkgcraft-go
# build and test
go test -v ./...
# benchmark
go test -bench=. ./...
```
