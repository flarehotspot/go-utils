package apiv1

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/flarehotspot/core/accounts"
// 	"github.com/flarehotspot/core/config"
// 	"github.com/flarehotspot/core/globals"
// 	"github.com/flarehotspot/core/sdk/api/http"
// 	"github.com/flarehotspot/core/utils/jsonwebtoken"
// 	"github.com/flarehotspot/core/web/middlewares"
// )

// func NewAdminAuthCtrl(g *globals.CoreGlobals) *AdminAuthCtrl {
// 	return &AdminAuthCtrl{g}
// }

// type AdminAuthCtrl struct {
// 	g *globals.CoreGlobals
// }

// func (c *AdminAuthCtrl) Login(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	username := r.PostFormValue("username")
// 	password := r.PostFormValue("password")
// 	acct, err := accounts.Find(username)
// 	if err != nil {
// 		c.ErrorUnauthorized(w, err.Error())
// 		return
// 	}

// 	if !acct.Auth(password) {
// 		err = errors.New(c.g.CoreApi.Utl.Translate("error", "invalid_login"))
// 		c.ErrorUnauthorized(w, err.Error())
// 		return
// 	}

// 	appcfg, err := config.ReadApplicationConfig()
// 	if err != nil {
// 		err = errors.New(c.g.CoreApi.Utl.Translate("error", "invalid_login"))
// 		c.ErrorUnauthorized(w, err.Error())
// 		return
// 	}

// 	payload := map[string]string{"username": username}
// 	token, err := jsonwebtoken.GenerateToken(payload, appcfg.Secret)
// 	if err != nil {
// 		c.ErrorUnauthorized(w, err.Error())
// 		return
// 	}

// 	sdkhttp.SetCookie(w, middlewares.AuthTokenCookie, token)
// 	data := map[string]string{"token": token}
// 	c.g.CoreApi.HttpApi().HttpResponse().Json(w, data, http.StatusOK)
// }

// func (c *AdminAuthCtrl) Logout(w http.ResponseWriter, r *http.Request) {
// 	sdkhttp.SetCookie(w, middlewares.AuthTokenCookie, "")
// 	data := map[string]string{"message": "Logout success"}
// 	c.g.CoreApi.HttpApi().HttpResponse().Json(w, data, http.StatusOK)
// }

// func (c *AdminAuthCtrl) IsAuthenticated(w http.ResponseWriter, r *http.Request) {
// 	data := map[string]string{"message": "Success"}
// 	c.g.CoreApi.HttpApi().HttpResponse().Json(w, data, http.StatusOK)
// }

// func (c *AdminAuthCtrl) ErrorUnauthorized(w http.ResponseWriter, msg string) {
// 	data := map[string]string{"error": msg}
// 	c.g.CoreApi.HttpApi().HttpResponse().Json(w, data, http.StatusUnauthorized)
// }
