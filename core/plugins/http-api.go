package plugins

import (
	"net/http"

	"github.com/flarehotspot/core/connmgr"
	"github.com/flarehotspot/core/db"
	"github.com/flarehotspot/core/db/models"
	"github.com/flarehotspot/core/payments"
	sdkhttp "github.com/flarehotspot/core/sdk/api/http"
	"github.com/gorilla/mux"
)

func NewHttpApi(api *PluginApi, db *db.Database, clnt *connmgr.ClientRegister, mdls *models.Models, dmgr *connmgr.ClientRegister, pmgr *payments.PaymentsMgr) *HttpApi {
	auth := NewAuthApi(api)
	httpRouter := NewHttpRouterApi(api, db, clnt)
	vueRouter := NewVueRouterApi(api)
	httpResp := NewHttpResponse(api)
	middlewares := NewPluginMiddlewares(api.db, mdls, dmgr, pmgr)
	return &HttpApi{
		api:         api,
		auth:        auth,
		httpRouter:  httpRouter,
		vueRouter:   vueRouter,
		httpResp:    httpResp,
		middlewares: middlewares,
	}
}

type HttpApi struct {
	api         *PluginApi
	auth        *AuthApi
	httpRouter  *HttpRouterApi
	vueRouter   *VueRouterApi
	httpResp    *HttpResponse
	middlewares *PluginMiddlewares
}

func (self *HttpApi) Auth() sdkhttp.IAuthApi {
	return self.auth
}

func (self *HttpApi) HttpRouter() sdkhttp.IHttpRouterApi {
	return self.httpRouter
}

func (self *HttpApi) VueRouter() sdkhttp.IVueRouterApi {
	return self.vueRouter
}

func (self *HttpApi) Helpers() sdkhttp.IHelpers {
	return NewViewHelpers(self.api)
}

func (self *HttpApi) Middlewares() sdkhttp.Middlewares {
	return self.middlewares
}

func (self *HttpApi) HttpResponse() sdkhttp.IHttpResponse {
	return self.httpResp
}

func (self *HttpApi) VueResponse() sdkhttp.IVueResponse {
	return NewVueResponse(self.api.HttpAPI.vueRouter)
}

func (self *HttpApi) MuxVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func (self *HttpApi) GetAdminNavs(r *http.Request) []sdkhttp.AdminNavCategory {
	return self.api.PluginsMgr.Utils().GetAdminNavs(r)
}
