package plugins

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/flarehotspot/core/sdk/api/http"
)

func NewVueRouterApi(api *PluginApi) *VueRouterApi {
	return &VueRouterApi{
		api:          api,
		adminRoutes:  []*VueRouteComponent{},
		portalRoutes: []*VueRouteComponent{},
	}
}

type VueRouterApi struct {
	api          *PluginApi
	adminRoutes  []*VueRouteComponent
	portalRoutes []*VueRouteComponent
	adminNavsFn  sdkhttp.VueAdminNavsFunc
	portalNavsFn sdkhttp.VuePortalItemsFunc
}

func (self *VueRouterApi) SetAdminRoutes(routes []sdkhttp.VueAdminRoute) {
	if routes != nil {
		dataRouter := self.api.HttpAPI.httpRouter.adminRouter.mux
		for _, r := range routes {
			route := NewVueRouteComponent(self.api, r.RouteName, r.RoutePath, r.HandlerFunc, r.Component, nil, nil)

			if _, ok := self.FindAdminRoute(route.VueRouteName); ok {
				log.Println("Warning: Admin route name \"" + r.RouteName + "\" already exists in admin routes ")
			}

			route.MountRoute(dataRouter, r.Middlewares...)
			self.adminRoutes = append(self.adminRoutes, route)
		}
	}
}

func (self *VueRouterApi) SetPortalRoutes(routes []sdkhttp.VuePortalRoute) {
	if routes != nil {

		pluginRouter := self.api.HttpAPI.httpRouter.pluginRouter
		compRouter := pluginRouter.mux.PathPrefix("/vue-route/portal-components").Subrouter()
		dataRouter := pluginRouter.mux.PathPrefix("/vue-route/portal-data").Subrouter()
		for _, r := range routes {
			route := NewVueRouteComponent(self.api, r.RouteName, r.RoutePath, r.HandlerFn, r.Component, nil, nil)

			if _, ok := self.FindPortalRoute(route.VueRouteName); ok {
				log.Println("Warning: Portal route name \"" + r.RouteName + "\" already exists in portal routes ")
			}

			compRouter.
				HandleFunc(route.HttpComponentPath, route.GetComponentHandler()).
				Methods("GET").
				Name(string(route.MuxCompRouteName))

			handler := http.HandlerFunc(route.GetDataHandler())

			var handlerFunc http.Handler
			if r.Middlewares != nil {
				for _, m := range r.Middlewares {
					handlerFunc = m(handlerFunc)
				}
			} else {
				handlerFunc = handler
			}

			dataRouter.
				Handle(route.HttpDataPath, handlerFunc).
				Methods("GET").
				Name(string(route.MuxDataRouteName))

			router := compRouter
			compR := router.Get(string(route.MuxCompRouteName))
			dataR := router.Get(string(route.MuxDataRouteName))
			comppath, _ := compR.GetPathTemplate()
			datapath, _ := dataR.GetPathTemplate()
			route.HttpComponentFullPath = comppath
			route.HttpDataFullPath = datapath
			self.portalRoutes = append(self.portalRoutes, route)
		}

	}
}

func (self *VueRouterApi) GetAdminRoutes() []*VueRouteComponent {
	return self.adminRoutes
}

func (self *VueRouterApi) GetPortalRoutes() []*VueRouteComponent {
	return self.portalRoutes
}

func (self *VueRouterApi) FindAdminRoute(routename string) (*VueRouteComponent, bool) {
	vueR := self.api.HttpAPI.vueRouter
	routeName := vueR.VueRouteName(routename)
	for _, route := range self.GetAdminRoutes() {
		if route.VueRouteName == routeName {
			return route, true
		}
	}
	return nil, false
}

func (self *VueRouterApi) AdminNavsFunc(fn sdkhttp.VueAdminNavsFunc) {
	self.adminNavsFn = fn
}

func (self *VueRouterApi) GetAdminNavs(r *http.Request) []sdkhttp.AdminNavItem {
	navs := []sdkhttp.AdminNavItem{}
	if self.adminNavsFn != nil {
		for _, nav := range self.adminNavsFn(r) {
			adminNav, ok := NewVueAdminNav(self.api, r, nav)
			if ok {
				navs = append(navs, adminNav)
			}
		}
	}

	return navs
}

func (self *VueRouterApi) FindPortalRoute(name string) (*VueRouteComponent, bool) {
	vueR := self.api.HttpAPI.vueRouter
	routeName := vueR.VueRouteName(name)

	for _, route := range self.GetPortalRoutes() {
		if route.VueRouteName == routeName {
			return route, true
		}
	}

	return nil, false
}

func (self *VueRouterApi) FindVueComponent(name string) (VueRouteComponent, bool) {
	return VueRouteComponent{}, true
}

func (self *VueRouterApi) PortalItemsFunc(fn sdkhttp.VuePortalItemsFunc) {
	self.portalNavsFn = fn
}

func (self *VueRouterApi) GetPortalItems(r *http.Request) []VuePortalItem {
	navs := []VuePortalItem{}

	if self.portalNavsFn != nil {
		for _, nav := range self.portalNavsFn(r) {
			navs = append(navs, NewVuePortalItem(self.api, r, nav))
		}

		return navs
	}

	return navs
}

func (self *VueRouterApi) FindVueRoute(name string) (*VueRouteComponent, bool) {
	vueR := self.api.HttpAPI.vueRouter
	routeName := vueR.VueRouteName(name)
	for _, route := range self.GetAdminRoutes() {
		if route.VueRouteName == routeName {
			return route, true
		}
	}
	for _, route := range self.GetPortalRoutes() {
		if route.VueRouteName == routeName {
			return route, true
		}
	}
	return nil, false
}

func (self *VueRouterApi) VueRouteName(name string) string {
	name = fmt.Sprintf("%s.%s", self.api.Pkg(), name)
	return name
}

func (self *VueRouterApi) VueRoutePath(path string) string {
	path = filepath.Join("/", self.api.Pkg(), path)
	return strings.TrimSuffix(path, "/")
}
func (self *VueRouterApi) HttpDataPath(path string) string {
	path = self.MuxPathFromVue(path)
	return strings.TrimSuffix(path, "/")
}

func (self *VueRouterApi) HttpComponentPath(path string) string {
	path = filepath.Join("/unwrapped/", path)
	if !strings.HasSuffix(path, ".vue") {
		path = path + ".vue"
	}
	return path
}

func (self *VueRouterApi) HttpWrapperRouteName(name string) string {
	name = fmt.Sprintf("%s.%s.%s", self.api.Pkg(), "wrapper", name, )
	return name
}

func (self *VueRouterApi) HttpWrapperRoutePath(name string) string {
	name = filepath.Join("/wrapper", name)
	if !strings.HasSuffix(name, ".vue") {
		name = name + ".vue"
	}
	return name
}

func (self *VueRouterApi) MuxPathFromVue(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	parts := strings.Split(path, "/")
	routPath := strings.Builder{}
	for _, p := range parts {
		if strings.HasPrefix(p, ":") {
			routPath.WriteString(fmt.Sprintf("{%s}", strings.TrimPrefix(p, ":")))
		} else {
			routPath.WriteString(p)
		}
		routPath.WriteString("/")
	}
	path = routPath.String()
	return strings.TrimSuffix(path, "/")
}
