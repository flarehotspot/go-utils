package web

import (
	"fmt"
	"net/http"

	"github.com/flarehotspot/core/plugins"
	"github.com/flarehotspot/core/web/router"
	"github.com/flarehotspot/core/web/routes"
	"github.com/gorilla/mux"
)

func SetupBootRoutes(g *plugins.CoreGlobals) {
	routes.BootRoutes(g)
	routes.CoreAssets(g)
}

func SetupAllRoutes(g *plugins.CoreGlobals) {
	routes.IndexRoutes(g)
	routes.AssetsRoutes(g)
	routes.PaymentRoutes(g)

	router.RootRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	router.RootRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		// met, err2 := route.GetMethods()
		fmt.Println(tpl)
		return nil
	})
}
