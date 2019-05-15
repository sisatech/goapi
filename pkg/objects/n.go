package objects

// NetworkInterface ..
type NetworkInterface struct {
	DisableTCPSegmentationOffload bool     `json:"disableTCPSegmentationOffload"`
	Gateway                       string   `json:"gateway"`
	HTTP                          []string `json:"http"`
	HTTPS                         []string `json:"https"`
	IP                            string   `json:"ip"`
	Mask                          string   `json:"mask"`
	MTU                           int      `json:"mtu"`
	TCP                           []string `json:"tcp"`
	UDP                           []string `json:"udp"`
}

// Node ..
type Node struct {
	Host string `json:"host"`
	Name string `json:"name"`
	Type string `json:"type"`
}
