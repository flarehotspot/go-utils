package apictrl

// import (
// 	"net/http"

// 	"github.com/flarehotspot/core/internal/plugins"
// 	"github.com/flarehotspot/core/internal/web/helpers"
// 	"github.com/flarehotspot/core/internal/web/response"
// 	models "github.com/flarehotspot/sdk/api/models"
// )

// type ClientApiCtrl struct {
// 	g *plugins.CoreGlobals
// }

// func NewClientApiCtrl(g *plugins.CoreGlobals) *ClientApiCtrl {
// 	return &ClientApiCtrl{g}
// }

// func (self *ClientApiCtrl) ClientData(w http.ResponseWriter, r *http.Request) {
// 	clnt, err := helpers.CurrentClient(r)
// 	if err != nil {
// 		ErrJson(w, r, err)
// 		return
// 	}

// 	dev, err := self.g.Models.Device().Find(r.Context(), clnt.Id())
// 	if err != nil {
// 		ErrJson(w, r, err)
// 		return
// 	}

// 	connected := self.g.ClientMgr.IsConnected(clnt)
// 	status := map[string]any{"connected": connected}
// 	devmap := models.DeviceToMap(dev)
// 	data := map[string]any{
// 		"device": devmap,
// 		"status": status,
// 	}
// 	response.Json(w, data, 200)
// }
