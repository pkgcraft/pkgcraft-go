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
	config *C.Config
	// cached fields
	_repos *Repos
}

// Return a new config for the system.
func NewConfig() (*Config, error) {
	ptr := C.pkgcraft_config_new()
	if ptr != nil {
		config := &Config{config: ptr, _repos: nil}
		runtime.SetFinalizer(config, func(c *Config) { C.pkgcraft_config_free(c.config) })
		return config, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

// Return the config's repo mapping.
func (c *Config) repos() (*Repos) {
	if c._repos == nil {
		c._repos = repos_from_config(c.config)
	}
	return c._repos
}

// Load repos from a portage-compatible repos.conf directory or file.
func (c *Config) load_repos_conf(path string) (map[string]Repo, error) {
	var length C.size_t

	path_str := C.CString(path)
	defer C.free(unsafe.Pointer(path_str))
	repos := C.pkgcraft_config_load_repos_conf(c.config, path_str, &length)

	if repos != nil {
		// force config repos refresh
		c._repos = nil

		m := repos_to_map(unsafe.Slice(repos, length))
		defer C.pkgcraft_repos_free(repos, length)
		return m, nil
	} else {
		s := C.pkgcraft_last_error()
		defer C.pkgcraft_str_free(s)
		return nil, errors.New(C.GoString(s))
	}
}

type Repos struct {
	config *C.Config
	// cached fields
	_repos map[string]Repo
}

// Return a Repos object for a given config.
func repos_from_config(config *C.Config) (*Repos) {
	var length C.size_t
	repos := C.pkgcraft_config_repos(config, &length)
	m := repos_to_map(unsafe.Slice(repos, length))
	defer C.pkgcraft_repos_free(repos, length)
	return &Repos{config: config, _repos: m}
}

// Convert an array of Repo pointers to a mapping.
func repos_to_map(repos []*C.Repo) map[string]Repo {
	m := make(map[string]Repo)
	for _, r := range repos {
		s := C.pkgcraft_repo_id(r)
		id := C.GoString(s)
		defer C.pkgcraft_str_free(s)
		m[id] = repo_from_ptr(r)
	}
	return m
}

// Return a Repos object length.
func (r *Repos) len() int {
	return len(r._repos)
}
