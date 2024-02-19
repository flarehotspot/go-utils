package web

import (
	"net/http"

	"github.com/flarehotspot/core/internal/plugins"
	"github.com/flarehotspot/core/internal/web/router"
	"github.com/flarehotspot/core/internal/web/routes"
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

	// router.RootRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// 	tpl, _ := route.GetPathTemplate()
	// 	// met, err2 := route.GetMethods()
	// 	fmt.Println(tpl)
	// 	return nil
	// })
}