package plugins

import (
	"net/http"
	"path/filepath"

	"github.com/flarehotspot/core/web/response"
)

const (
	rootjson = "$$response"
)

func NewVueResponse(vr *VueRouterApi) *VueResponse {
	return &VueResponse{vr, map[string]any{
		rootjson: map[string]any{},
	}}
}

type VueResponse struct {
	router *VueRouterApi
	data   map[string]any
}

func (self *VueResponse) FlashMsg(msgType string, msg string) {
	newdata := self.data[rootjson].(map[string]any)
	newdata["flash"] = map[string]string{
		"type": msgType, // "success", "error", "warning", "info
		"msg":  msg,
	}
	self.data[rootjson] = newdata
}

func (self *VueResponse) Json(w http.ResponseWriter, data any, status int) {
	if data == nil {
		data = map[string]any{}
	}
	newdata := self.data[rootjson].(map[string]any)
	newdata["data"] = data
	data = map[string]any{
		rootjson: newdata,
	}
	response.Json(w, data, status)
}

func (self *VueResponse) Component(w http.ResponseWriter, vuefile string, data any) {
	vuefile = self.router.api.Utl.Resource(filepath.Join("components", vuefile))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	helpers := self.router.api.HttpAPI.Helpers()
	response.Text(w, vuefile, helpers, data)
}

func (res *VueResponse) Redirect(w http.ResponseWriter, routename string, pairs ...string) {
	route, ok := res.router.FindVueRoute(routename)
	if !ok {
		response.ErrorJson(w, "Vue route \""+routename+"\" not found", 500)
		return
	}

	params := map[string]string{}
	for i := 0; i < len(pairs); i += 2 {
		params[pairs[i]] = pairs[i+1]
	}

	newdata := res.data[rootjson].(map[string]any)
	newdata["redirect"] = true
	newdata["route_name"] = route.VueRouteName
	newdata["params"] = params
	data := map[string]any{rootjson: newdata}

	response.Json(w, data, http.StatusOK)
}