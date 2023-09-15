package network

import (
	"log"
	"sync"
	"time"

	"github.com/flarehotspot/core/utils/nftables"
	"github.com/flarehotspot/core/sdk/api/network"
)

type TrafficMgr struct {
	mu        sync.RWMutex
	ticker    *time.Ticker
	listners  []chan network.TrafficData
	prevStats *nftables.StatResult
}

func NewTrafficMgr() *TrafficMgr {
	return &TrafficMgr{}
}

func (self *TrafficMgr) Start() {
	go func() {
		self.mu.Lock()
		self.ticker = time.NewTicker(5 * time.Second)
		self.mu.Unlock()

		for range self.ticker.C {
			go self.MakeTrafficData()
		}
	}()
}

func (self *TrafficMgr) Listen() <-chan network.TrafficData {
	retCh := make(chan chan network.TrafficData)
	go func() {
		self.mu.Lock()
		defer self.mu.Unlock()
		ch := make(chan network.TrafficData)
		self.listners = append(self.listners, ch)
		retCh <- ch
	}()

	return <-retCh
}

func (self *TrafficMgr) MakeTrafficData() {
	self.mu.Lock()
	defer self.mu.Unlock()

	if len(self.listners) == 0 {
		return
	}

	stats, err := nftables.GetStats()
	if err != nil {
		log.Println(err)
		return
	}

	prevStats := &nftables.StatResult{
		MacStats: make(map[string]nftables.StatData),
		IpStats:  make(map[string]nftables.StatData),
	}

	if self.prevStats != nil {
		prevStats = self.prevStats
	}

	trfc := network.TrafficData{
		Download: make(map[string]network.ClientStat),
		Upload:   make(map[string]network.ClientStat),
	}

	for mac, stat := range stats.MacStats {
		prev, ok := prevStats.MacStats[mac]
		if ok {
			// If current stat is less than prev, user may have been reconnected.
			// In this case we discard previous stats.
			if stat.Packets < prev.Packets || stat.Bytes < prev.Bytes {
				trfc.Upload[mac] = network.ClientStat{Packets: stat.Packets, Bytes: stat.Bytes}
				continue
			}

			pkts := stat.Packets - prev.Packets
			byts := stat.Bytes - prev.Bytes
			trfc.Upload[mac] = network.ClientStat{Packets: pkts, Bytes: byts}
		} else {
			trfc.Upload[mac] = network.ClientStat{Packets: stat.Packets, Bytes: stat.Bytes}
		}
	}

	for ip, stat := range stats.IpStats {
		prev, ok := prevStats.IpStats[ip]
		if ok {
			// If current stat is less than prev, user may have been reconnected.
			// In this case we discard previous stats.
			if stat.Packets < prev.Packets || stat.Bytes < prev.Bytes {
				trfc.Download[ip] = network.ClientStat{Packets: stat.Packets, Bytes: stat.Bytes}
				continue
			}

			pkts := stat.Packets - prev.Packets
			byts := stat.Bytes - prev.Bytes
			trfc.Download[ip] = network.ClientStat{Packets: pkts, Bytes: byts}
		} else {
			trfc.Download[ip] = network.ClientStat{Packets: stat.Packets, Bytes: stat.Bytes}
		}
	}

	for _, ch := range self.listners {
		ch <- trfc
	}

	self.prevStats = &stats
}

// func (self *DataConnMgr) nftStatToMap