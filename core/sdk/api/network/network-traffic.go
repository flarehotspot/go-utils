package network

type ClientStat struct {
	Packets uint
	Bytes   uint
}

// TrafficData represents a network traffic data.
type TrafficData struct {

	// Download is a map of client IP addresses to download statistics.
	Download map[string]ClientStat

	// Upload is a map of client MAC addresses to upload statistics.
	Upload map[string]ClientStat
}

// ITrafficApi is the interface for the network traffic API.
// It can be used to listen to network traffic.
// It emits network traffic data every 5 seconds.
type ITrafficApi interface {
	Listen() <-chan TrafficData
}