package pkgcraft

// #cgo pkg-config: pkgcraft
// #include <pkgcraft.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type Config struct {
	ptr *C.Config
	// cached fields
	_repos *Repos
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

// Return the config's repo mapping.
func (c *Config) Repos() *Repos {
	if c._repos == nil {
		c._repos = repos_from_config(c)
	}
	return c._repos
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
	c._repos = nil

	return repo_from_ptr(ptr, false), nil
}

// Load repos from a portage-compatible repos.conf directory or file.
func (c *Config) LoadReposConf(path string) (map[string]Repo, error) {
	var length C.size_t

	path_str := C.CString(path)
	defer C.free(unsafe.Pointer(path_str))
	repos := C.pkgcraft_config_load_repos_conf(c.ptr, path_str, &length)

	if repos != nil {
		// force config repos refresh
		c._repos = nil

		m := repos_to_map(unsafe.Slice(repos, length), false)
		defer C.pkgcraft_repos_free(repos, length)
		return m, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

type Repos struct {
	config *Config
	// cached fields
	_repos map[string]Repo
}

// Return a Repos object for a given config.
func repos_from_config(config *Config) *Repos {
	var length C.size_t
	repos := C.pkgcraft_config_repos(config.ptr, &length)
	m := repos_to_map(unsafe.Slice(repos, length), true)
	defer C.pkgcraft_repos_free(repos, length)
	return &Repos{config: config, _repos: m}
}

// Convert an array of Repo pointers to a mapping.
func repos_to_map(repos []*C.Repo, ref bool) map[string]Repo {
	m := make(map[string]Repo)
	for _, r := range repos {
		s := C.pkgcraft_repo_id(r)
		id := C.GoString(s)
		defer C.pkgcraft_str_free(s)
		m[id] = repo_from_ptr(r, ref)
	}
	return m
}

// Return a Repos object length.
func (r *Repos) Len() int {
	return len(r._repos)
}
