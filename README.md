# pkgcraft-go

Go bindings for pkgcraft.

To building the bindings, run the following commands:

```bash
git clone --recurse-submodules https://github.com/pkgcraft/scallop.git
git clone https://github.com/pkgcraft/pkgcraft.git
git clone https://github.com/pkgcraft/pkgcraft-c.git
git clone https://github.com/pkgcraft/pkgcraft-go.git

# install cargo-c
cargo install cargo-c

# build pkgcraft-c library
cd pkgcraft-go
cargo cinstall --prefix="${PWD}/pkgcraft" --pkgconfigdir="${PWD}/pkgcraft" --manifest-path=../pkgcraft-c/Cargo.toml
export PKG_CONFIG_PATH="${PWD}/pkgcraft"
export LD_LIBRARY_PATH="${PWD}/pkgcraft/lib"

# build go bindings
go build
```
