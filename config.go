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
	ptr   *C.Config
	Repos map[string]*BaseRepo
	ReposEbuild map[string]*EbuildRepo
	ReposFake map[string]*FakeRepo
}

// Return a new config for the system.
func NewConfig() (*Config, error) {
	ptr := C.pkgcraft_config_new()
	if ptr != nil {
		config := &Config{ptr: ptr}
		_, file, line, _ := runtime.Caller(1)
		runtime.SetFinalizer(config, func(config *Config) {
			panic(fmt.Sprintf("%s:%d: unclosed config object", file, line))
		})
		return config, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

// Free a config object's encapsulated C pointer.
func (config *Config) Close() {
	C.pkgcraft_config_free(config.ptr)
	runtime.SetFinalizer(config, nil)
}

// Add an external repo via its file path.
func (config *Config) AddRepoPath(path string, id string, priority int) error {
	path_str := C.CString(path)
	defer C.free(unsafe.Pointer(path_str))
	id_str := C.CString(id)
	defer C.free(unsafe.Pointer(id_str))

	ptr := C.pkgcraft_config_add_repo_path(config.ptr, id_str, C.int(priority), path_str)
	if ptr == nil {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return errors.New(C.GoString(s))
	}

	config.updateRepos()
	return nil
}

// Load repos from a portage-compatible repos.conf directory or file.
func (config *Config) LoadReposConf(path string) error {
	var length C.size_t

	path_str := C.CString(path)
	defer C.free(unsafe.Pointer(path_str))
	repos := C.pkgcraft_config_load_repos_conf(config.ptr, path_str, &length)

	if repos != nil {
		config.updateRepos()
		C.pkgcraft_repos_free(repos, length)
		return nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return errors.New(C.GoString(s))
	}
}

// Update the repo maps for a config.
func (config *Config) updateRepos() {
	var length C.size_t
	c_repos := C.pkgcraft_config_repos(config.ptr, &length)
	config.Repos = repos_to_map(unsafe.Slice(c_repos, length))
	config.ReposEbuild = make(map[string]*EbuildRepo)
	config.ReposFake = make(map[string]*FakeRepo)
	for id, r := range config.Repos {
		switch format := r.format; format {
			case RepoFormatEbuild: config.ReposEbuild[id] = &EbuildRepo{r}
			case RepoFormatFake: config.ReposFake[id] = &FakeRepo{r}
		}
	}
	C.pkgcraft_repos_free(c_repos, length)
}

// Convert an array of Repo pointers to a mapping.
func repos_to_map(repos []*C.Repo) map[string]*BaseRepo {
	m := make(map[string]*BaseRepo)
	for _, r := range repos {
		s := C.pkgcraft_repo_id(r)
		id := C.GoString(s)
		defer C.pkgcraft_str_free(s)
		m[id] = repo_from_ptr(r)
	}
	return m
}
