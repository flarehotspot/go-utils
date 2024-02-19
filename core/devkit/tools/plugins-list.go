package tools

import (
	"path/filepath"

	sdkfs "github.com/flarehotspot/sdk/utils/fs"
)

func PluginPathList() []string {
	searchPaths := []string{"system", "plugins"}
	pluginPaths := []string{}

	for _, sp := range searchPaths {
		var dirs []string
		if err := sdkfs.LsDirs(sp, &dirs, false); err != nil {
			continue
		}

		for _, dir := range dirs {
			pluginJson := filepath.Join(dir, "plugin.json")
			modFile := filepath.Join(dir, "go.mod")

			if sdkfs.Exists(pluginJson) && sdkfs.Exists(modFile) {
				pluginPaths = append(pluginPaths, dir)
			}
		}
	}

	return pluginPaths
}
