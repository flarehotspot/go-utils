package portalctrl

import (
	"net/http"

	"github.com/flarehotspot/core/config/themecfg"
	"github.com/flarehotspot/core/globals"
	"github.com/flarehotspot/core/plugins"
)

func NewPortalCtrl(g *globals.CoreGlobals) PortalCtrl {
	return PortalCtrl{g}
}

type PortalCtrl struct {
	g *globals.CoreGlobals
}

func (c *PortalCtrl) IndexPage(w http.ResponseWriter, r *http.Request) {
	themePkg := themecfg.Read().Portal
	themePlugin := c.g.PluginMgr.FindByPkg(themePkg)
	themesApi := themePlugin.ThemesApi().(*plugins.ThemesApi)
	portalComponent, ok := themesApi.GetPortalComponent()
	if !ok {
		http.Error(w, "No portal theme component path defined", 500)
		return
	}

	scripts := []string{}
	styles := []string{}

	if portalComponent.ThemeAssets != nil {
		if portalComponent.ThemeAssets.Scripts != nil {
			for _, script := range *portalComponent.ThemeAssets.Scripts {
				jsPath := themePlugin.HttpApi().Helpers(w, r).AssetPath(script)
				scripts = append(scripts, jsPath)
			}
		}

		if portalComponent.ThemeAssets.Styles != nil {
			for _, style := range *portalComponent.ThemeAssets.Styles {
				cssPath := themePlugin.HttpApi().Helpers(w, r).AssetPath(style)
				styles = append(styles, cssPath)
			}
		}
	}

	vdata := map[string]any{
		"CoreApi":      c.g.CoreApi,
		"ThemeScripts": scripts,
		"ThemeStyles":  styles,
	}

	api := c.g.CoreApi
	api.HttpApi().Respond().View(w, r, "portal/layout.html", vdata)
}
