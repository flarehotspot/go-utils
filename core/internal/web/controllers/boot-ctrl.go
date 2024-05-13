package controllers

import (
	"net/http"

	"github.com/flarehotspot/core/internal/plugins"
	sse "github.com/flarehotspot/core/internal/utils/sse"
	"github.com/flarehotspot/core/internal/web/helpers"
	"github.com/flarehotspot/core/internal/web/routes/urls"
)

func NewBootCtrl(g *plugins.CoreGlobals, pmgr *plugins.PluginsMgr, api *plugins.PluginApi) BootCtrl {
	return BootCtrl{g.BootProgress, pmgr, api}
}

type BootCtrl struct {
	bp   *plugins.BootProgress
	pmgr *plugins.PluginsMgr
	api  *plugins.PluginApi
}

func (ctrl *BootCtrl) IndexPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"title":  "Booting",
		"status": ctrl.bp.Status(),
		"done":   ctrl.bp.IsDone(),
	}

	ctrl.api.Http().HttpResponse().View(w, r, "booting/index.html", data)
}

func (ctrl *BootCtrl) SseHandler(w http.ResponseWriter, r *http.Request) {
	s, err := sse.NewSocket(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctrl.bp.AddSocket(s)
	s.Listen()
}

func (ctrl *BootCtrl) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := ctrl.bp.IsDone()

		if r.Method == "GET" && !helpers.IsAssetPath(r.URL.Path) {
			if !done && r.URL.Path != urls.BOOT_URL && r.URL.Path != urls.BOOT_STATUS_URL {
				http.Redirect(w, r, urls.BOOT_URL, http.StatusSeeOther)
				return
			}

			if done && r.URL.Path == urls.BOOT_URL {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}