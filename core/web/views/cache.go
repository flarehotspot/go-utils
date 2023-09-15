package views

import (
	"errors"
	"path/filepath"
	"sync"

	"html/template"

	"github.com/flarehotspot/core/sdk/utils/paths"
)

var (
	vmap           = sync.Map{}
	cachedir       = filepath.Join(paths.CacheDir, "views")
	errNoViewCache = errors.New("View cache not available.")
)

func WriteViewCache(tmpl *template.Template, views ...*ViewInput) error {
	files := viewFiles(views...)
	hash, err := filesHash(files...)
	if err != nil {
		return err
	}

	vkey, err := viewsHash(views...)
	if err != nil {
		return err
	}

	vcache := &viewCache{tmpl, hash}
	vmap.Store(vkey, vcache)

	return nil
}

func GetViewCache(views ...*ViewInput) (*template.Template, error) {
	files := viewFiles(views...)
	hash, err := filesHash(files...)
	if err != nil {
		return nil, err
	}

	vkey, err := viewsHash(views...)
	if err != nil {
		return nil, err
	}

	sym, ok := vmap.Load(vkey)
	if !ok {
		return nil, errNoViewCache
	}

	vcache := sym.(*viewCache)
	if vcache.hash != hash {
		return nil, errNoViewCache
	}

	return vcache.tmpl, nil
}

func viewAssets(f string) []string {
	files := []string{}
	m := ViewManifest(f)

	for _, src := range m.Scripts {
		files = append(files, filepath.Join(f, "..", src))
	}

	for _, src := range m.Styles {
		files = append(files, filepath.Join(f, "..", src))
	}

	for _, folder := range m.Folders {
		files = append(files, filepath.Join(f, "..", folder.Src))
	}

	return files
}