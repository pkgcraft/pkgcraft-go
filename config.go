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
	ptr *C.Config
	// cached fields
	Repos map[string]*BaseRepo
}

// Return a new config for the system.
func NewConfig() (*Config, error) {
	ptr := C.pkgcraft_config_new()
	if ptr != nil {
		config := &Config{ptr: ptr}
		runtime.SetFinalizer(config, func(c *Config) { C.pkgcraft_config_free(c.ptr) })
		return config, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

// Add an external repo via its file path.
func (c *Config) AddRepoPath(path string, id string, priority int) (Repo, error) {
	path_str := C.CString(path)
	defer C.free(unsafe.Pointer(path_str))
	id_str := C.CString(id)
	defer C.free(unsafe.Pointer(id_str))

	ptr := C.pkgcraft_config_add_repo_path(c.ptr, id_str, C.int(priority), path_str)
	if ptr == nil {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}

	// force config repos refresh
	c.Repos = repos_from_config(c)

	return repo_from_ptr(ptr, false), nil
}

// Load repos from a portage-compatible repos.conf directory or file.
func (c *Config) LoadReposConf(path string) error {
	var length C.size_t

	path_str := C.CString(path)
	defer C.free(unsafe.Pointer(path_str))
	repos := C.pkgcraft_config_load_repos_conf(c.ptr, path_str, &length)

	if repos != nil {
		// force config repos refresh
		c.Repos = repos_from_config(c)
		C.pkgcraft_repos_free(repos, length)
		return nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return errors.New(C.GoString(s))
	}
}

// Return a configured ebuild repo from a given id.
func (c *Config) GetEbuildRepo(id string) (*EbuildRepo, error) {
	repo, exists := c.Repos[id]
	if exists {
		if repo.format == RepoFormatEbuild {
			return &EbuildRepo{repo}, nil
		} else {
			return nil, fmt.Errorf("invalid repo type: %s", id)
		}
	} else {
		return nil, fmt.Errorf("nonexistent repo: %s", id)
	}
}

// Return a configured fake repo from a given id.
func (c *Config) GetFakeRepo(id string) (*FakeRepo, error) {
	repo, exists := c.Repos[id]
	if exists {
		if repo.format == RepoFormatFake {
			return &FakeRepo{repo}, nil
		} else {
			return nil, fmt.Errorf("invalid repo type: %s", id)
		}
	} else {
		return nil, fmt.Errorf("nonexistent repo: %s", id)
	}
}

// Return a Repos object for a given config.
func repos_from_config(config *Config) map[string]*BaseRepo {
	var length C.size_t
	repos := C.pkgcraft_config_repos(config.ptr, &length)
	m := repos_to_map(unsafe.Slice(repos, length), true)
	C.pkgcraft_repos_free(repos, length)
	return m
}

// Convert an array of Repo pointers to a mapping.
func repos_to_map(repos []*C.Repo, ref bool) map[string]*BaseRepo {
	m := make(map[string]*BaseRepo)
	for _, r := range repos {
		s := C.pkgcraft_repo_id(r)
		id := C.GoString(s)
		defer C.pkgcraft_str_free(s)
		m[id] = repo_from_ptr(r, ref)
	}
	return m
}
