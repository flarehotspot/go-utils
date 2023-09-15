package plugins

import (
	"github.com/flarehotspot/core/utils/uci"
	ucisdk "github.com/flarehotspot/core/sdk/api/uci"
	gouci "github.com/flarehotspot/core/sdk/libs/go-uci"
)

type UciApi struct {
	networkApi  *uci.UciNetworkApi
	dhcpApi     *uci.UciDhcpApi
	wirelessApi *uci.UciWirelessApi
}

func NewUciApi() *UciApi {
	return &UciApi{}
}

func (self *UciApi) Network() ucisdk.INetworkApi {
	return self.networkApi
}

func (self *UciApi) Dhcp() ucisdk.IDhcpApi {
	return self.dhcpApi
}

func (self *UciApi) Wireless() ucisdk.IWirelessApi {
	return self.wirelessApi
}

func (self *UciApi) Uci() gouci.Tree {
	return uci.UciTree
}