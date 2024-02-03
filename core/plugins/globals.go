package plugins

import (
	"github.com/flarehotspot/core/connmgr"
	"github.com/flarehotspot/core/db"
	"github.com/flarehotspot/core/db/models"
	"github.com/flarehotspot/core/network"
	paths "github.com/flarehotspot/core/sdk/utils/paths"
)

type CoreGlobals struct {
	Db             *db.Database
	CoreAPI        *PluginApi
	ClientRegister *connmgr.ClientRegister
	ClientMgr      *connmgr.ClientMgr
	TrafficMgr     *network.TrafficMgr
	BootProgress   *BootProgress
	Models         *models.Models
	PluginMgr      *PluginsMgr
	PaymentsMgr    *PaymentsMgr
}

func New() *CoreGlobals {
	db, _ := db.NewDatabase()
	bp := NewBootProgress()
	mdls := models.New(db)
	clntReg := connmgr.NewClientRegister(db, mdls)
	clntMgr := connmgr.NewClientMgr(db, mdls)
	trfcMgr := network.NewTrafficMgr()
	pmtMgr := NewPaymentMgr()

	trfcMgr.Start()
	clntMgr.ListenTraffic(trfcMgr)

	plgnMgr := NewPluginMgr(db, mdls, pmtMgr, clntReg, clntMgr, trfcMgr)
	coreApi := NewPluginApi(paths.CoreDir, plgnMgr, trfcMgr)
	plgnMgr.InitCoreApi(coreApi)

	return &CoreGlobals{db, coreApi, clntReg, clntMgr, trfcMgr, bp, mdls, plgnMgr, pmtMgr}
}