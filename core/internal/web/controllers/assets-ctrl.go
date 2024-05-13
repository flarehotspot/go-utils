package controllers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/flarehotspot/core/internal/plugins"
	"github.com/flarehotspot/core/internal/web/response"
	"github.com/flarehotspot/sdk/utils/fs"
	"github.com/gorilla/mux"
)

func NewAssetsCtrl(g *plugins.CoreGlobals) *AssetsCtrl {
	return &AssetsCtrl{g}
}

type AssetsCtrl struct {
	g *plugins.CoreGlobals
}

func (ctrl *AssetsCtrl) GetFavicon(w http.ResponseWriter, r *http.Request) {
	contents, err := os.ReadFile(ctrl.g.CoreAPI.Utl.Resource("assets/images/default-favicon-32x32.png"))
	if err != nil {
		response.ErrorHtml(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(contents)
}

func (ctrl *AssetsCtrl) AssetWithHelpers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pkg := vars["pkg"]
	assetPath := vars["path"]
	pluginApi, ok := ctrl.g.PluginMgr.FindByPkg(pkg)
	if !ok {
		http.Error(w, "Plugin not found: "+pkg, 404)
		return
	}

	assetPath = filepath.Join(pluginApi.Resource("assets"), assetPath)
	if !sdkfs.Exists(assetPath) {
		http.Error(w, "Asset not found: "+assetPath, 404)
		return
	}

	response.File(w, assetPath, ctrl.g.CoreAPI.Http().Helpers(), nil)
}

func (ctrl *AssetsCtrl) VueComponent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pkg := vars["pkg"]
	componentPath := vars["path"]
	pluginApi, ok := ctrl.g.PluginMgr.FindByPkg(pkg)
	if !ok {
		ctrl.g.CoreAPI.HttpAPI.VueResponse().Component(w, "empty-component.vue", vars)
		return
	}

	res := pluginApi.Http().VueResponse()
	res.Component(w, componentPath, nil)
}