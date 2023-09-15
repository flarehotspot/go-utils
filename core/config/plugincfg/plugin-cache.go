package plugincfg

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/flarehotspot/core/sdk/utils/fs"
	"github.com/flarehotspot/core/sdk/utils/paths"
	"github.com/flarehotspot/core/sdk/utils/strings"
	jobque "github.com/flarehotspot/core/utils/job-que"
)

var (
	que       = jobque.NewJobQues()
	cacheJson = filepath.Join(paths.CacheDir, "plugins.json")
)

// map key format: [git|store]::[url]#[ref|version]
//
// Ex. git::https://github.com/user/pkg.git#ref
// Ex. store::com.flarego.my-plugin#1.0.0

type CacheInfo struct {
	Name    string `json:"name"`
	Package string `json:"package"`
	Version string `json:"version"`
}

func WriteCache(def *PluginSrcDef, info *PluginInfo) error {
	_, err := que.Exec(func() (interface{}, error) {
		cache := &CacheInfo{
			Name:    info.Name,
			Package: info.Package,
			Version: info.Version,
		}

		cfg, err := readCache()
		if err != nil {
			cfg = make(map[string]*CacheInfo)
		}

		key := getCacheKey(def)
		cfg[key] = cache

		b, err := json.Marshal(cfg)
		if err != nil {
			return nil, err
		}

		dir := filepath.Dir(cacheJson)
		if !fs.Exists(dir) {
			os.MkdirAll(dir, 0755)
		}

		err = os.WriteFile(cacheJson, b, 0644)
		return nil, err
	})

	return err
}

func GetCacheInfo(def *PluginSrcDef) (*CacheInfo, bool) {
	sym, err := que.Exec(func() (interface{}, error) {
		cfg, err := readCache()
		if err != nil {
			return nil, err
		}

		info, ok := cfg[getCacheKey(def)]
		if !ok {
			return nil, errors.New("plugin cache info not found")
		}

		return info, nil
	})

	if err != nil {
		return nil, false
	}

	return sym.(*CacheInfo), true
}

func getCacheKey(def *PluginSrcDef) string {
	var key string
	if def.Src == PluginSrcGit {
		key = "git::" + def.GitURL + "#" + def.GitRef
	}

	if def.Src == PluginSrcStore {
		key = "store::" + def.StorePackage + "#" + def.StoreVersion
	}

	return strings.Sha1Hash(key)
}

func readCache() (map[string]*CacheInfo, error) {
	b, err := os.ReadFile(cacheJson)
	if err != nil {
		return nil, err
	}

	var caches map[string]*CacheInfo
	err = json.Unmarshal(b, &caches)
	if err != nil {
		return nil, err
	}

	return caches, nil
}