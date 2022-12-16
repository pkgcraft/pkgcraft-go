package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

type Config struct {
	ptr         *C.Config
	Repos       map[string]*BaseRepo
	ReposEbuild map[string]*EbuildRepo
	ReposFake   map[string]*FakeRepo
}

// Return a new config for the system.
func NewConfig() *Config {
	ptr := C.pkgcraft_config_new()
	config := &Config{ptr: ptr}

	// force caller to explicitly close Config object, otherwise a panic occurs
	_, file, line, _ := runtime.Caller(1)
	runtime.SetFinalizer(config, func(self *Config) {
		panic(fmt.Sprintf("%s:%d: unclosed config object", file, line))
	})

	return config
}

// Free a config object's encapsulated C pointer.
func (self *Config) Close() {
	C.pkgcraft_config_free(self.ptr)
	runtime.SetFinalizer(self, nil)
}

// Add an external repo via its file path.
func (self *Config) AddRepoPath(path string, id string, priority int) error {
	path_str := C.CString(path)
	defer C.free(unsafe.Pointer(path_str))
	id_str := C.CString(id)
	defer C.free(unsafe.Pointer(id_str))

	ptr := C.pkgcraft_config_add_repo_path(self.ptr, id_str, C.int(priority), path_str)
	if ptr == nil {
		err := C.pkgcraft_error_last()
		defer C.pkgcraft_error_free(err)
		return errors.New(C.GoString(err.message))
	}

	self.updateRepos()
	return nil
}

// Add an external repo.
func (self *Config) AddRepo(repo repoPtr) error {
	ptr := C.pkgcraft_config_add_repo(self.ptr, repo.p())
	if ptr == nil {
		err := C.pkgcraft_error_last()
		defer C.pkgcraft_error_free(err)
		return errors.New(C.GoString(err.message))
	}

	self.updateRepos()
	return nil
}

// Load repos from a portage-compatible repos.conf directory or file.
func (self *Config) LoadReposConf(path string) error {
	var length C.size_t

	path_str := C.CString(path)
	defer C.free(unsafe.Pointer(path_str))
	repos := C.pkgcraft_config_load_repos_conf(self.ptr, path_str, &length)

	if repos != nil {
		self.updateRepos()
		C.pkgcraft_repos_free(repos, length)
		return nil
	} else {
		err := C.pkgcraft_error_last()
		defer C.pkgcraft_error_free(err)
		return errors.New(C.GoString(err.message))
	}
}

// Update the repo maps for a config.
func (self *Config) updateRepos() {
	var length C.size_t
	c_repos := C.pkgcraft_config_repos(self.ptr, &length)
	repos := make(map[string]*BaseRepo)
	for _, r := range unsafe.Slice(c_repos, length) {
		s := C.pkgcraft_repo_id(r)
		id := C.GoString(s)
		defer C.pkgcraft_str_free(s)
		repos[id] = repoFromPtr(r)
	}
	C.pkgcraft_repos_free(c_repos, length)

	repos_ebuild := make(map[string]*EbuildRepo)
	repos_fake := make(map[string]*FakeRepo)
	for id, r := range repos {
		switch format := r.format; format {
		case RepoFormatEbuild:
			repos_ebuild[id] = &EbuildRepo{r}
		case RepoFormatFake:
			repos_fake[id] = &FakeRepo{r}
		default:
			panic(fmt.Sprintf("unknown repo format: %d", format))
		}
	}

	self.Repos = repos
	self.ReposEbuild = repos_ebuild
	self.ReposFake = repos_fake
}
