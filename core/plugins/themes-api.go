package plugins

import (
	themes "github.com/flarehotspot/core/sdk/api/themes"
)

func NewThemesApi(api *PluginApi) *ThemesApi {
	return &ThemesApi{api: api}
}

type ThemesApi struct {
	api         *PluginApi
	AdminTheme  themes.AdminTheme
	PortalTheme themes.PortalTheme

	AdminLayoutRoute    *VueRouteComponent
	AdminDashboardRoute *VueRouteComponent
	AdminLoginRoute     *VueRouteComponent

	PortalLayoutRoute *VueRouteComponent
	PortalIndexRoute  *VueRouteComponent
}

func (t *ThemesApi) NewAdminTheme(theme themes.AdminTheme) {
	adminRouter := t.api.HttpAPI.httpRouter.adminRouter.mux
	compRouter := t.api.HttpAPI.httpRouter.pluginRouter.mux

	layoutComp := NewVueRouteComponent(t.api, theme.LayoutComponent.RouteName, "/theme/layout", theme.LayoutComponent.HandlerFunc, theme.LayoutComponent.ComponentPath, nil, nil)
	layoutComp.MountRoute(compRouter)

	loginComp := NewVueRouteComponent(t.api, theme.LoginComponent.RouteName, "/theme/login", theme.LoginComponent.HandlerFunc, theme.LoginComponent.ComponentPath, nil, nil)
	loginComp.MountRoute(compRouter)

	dashComp := NewVueRouteComponent(t.api, theme.DashboardComponent.RouteName, "/theme/dashboard", theme.DashboardComponent.HandlerFunc, theme.DashboardComponent.ComponentPath, nil, nil)
	dashComp.MountRoute(adminRouter)

	t.AdminLayoutRoute = layoutComp
	t.AdminDashboardRoute = dashComp
	t.AdminLoginRoute = loginComp
	t.api.HttpAPI.vueRouter.AddAdminRoutes(dashComp)
	t.api.HttpAPI.vueRouter.SetLoginRoute(loginComp)
	t.AdminTheme = theme
}

func (t *ThemesApi) NewPortalTheme(theme themes.PortalTheme) {
	compRouter := t.api.HttpAPI.httpRouter.pluginRouter.mux.PathPrefix("/portal/vue/components").Subrouter()

	layoutComp := NewVueRouteComponent(t.api, theme.LayoutComponent.RouteName, "/theme/layout", theme.LayoutComponent.HandlerFunc, theme.LayoutComponent.ComponentPath, nil, nil)
	layoutComp.MountRoute(compRouter)

	indexComp := NewVueRouteComponent(t.api, theme.IndexComponent.RouteName, "/theme/index", theme.IndexComponent.HandlerFunc, theme.IndexComponent.ComponentPath, nil, nil)
	indexComp.MountRoute(compRouter)

	t.PortalLayoutRoute = layoutComp
	t.PortalIndexRoute = indexComp
	t.api.HttpAPI.vueRouter.AddPortalRoutes(indexComp)
	t.PortalTheme = theme
}

func (t *ThemesApi) GetAdminThemeAssets() themes.ThemeAssets {
	assets := themes.ThemeAssets{Scripts: []string{}, Styles: []string{}}
	if t.AdminTheme.ThemeAssets != nil {
		if t.AdminTheme.ThemeAssets.Scripts != nil {
			assets.Scripts = t.AdminTheme.ThemeAssets.Scripts
		}
		if t.AdminTheme.ThemeAssets.Styles != nil {
			assets.Styles = t.AdminTheme.ThemeAssets.Styles
		}
	}
	return assets
}

func (t *ThemesApi) GetPortalThemeAssets() themes.ThemeAssets {
	assets := themes.ThemeAssets{Scripts: []string{}, Styles: []string{}}
	if t.PortalTheme.ThemeAssets != nil {
		if t.PortalTheme.ThemeAssets.Scripts != nil {
			assets.Scripts = t.PortalTheme.ThemeAssets.Scripts
		}
		if t.PortalTheme.ThemeAssets.Styles != nil {
			assets.Styles = t.PortalTheme.ThemeAssets.Styles
		}
	}
	return assets
}
