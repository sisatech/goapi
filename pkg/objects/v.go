package objects

// VM ..
type VM struct {
	Args      string       `json:"args"`
	Author    string       `json:"author"`
	Binary    string       `json:"binary"`
	CPUs      int          `json:"cpus"`
	Created   string       `json:"created"`
	Date      string       `json:"date"`
	Disk      string       `json:"disk"`
	Download  string       `json:"download"`
	Env       []string     `json:"env"`
	Hostname  string       `json:"hostname"`
	ID        string       `json:"id"`
	Instance  string       `json:"instance"`
	Kernel    string       `json:"kernel"`
	LogFile   string       `json:"logFile"`
	Name      string       `json:"name"`
	Networks  []VMNetwork  `json:"vmNetwork"`
	Platform  string       `json:"platform"`
	RAM       string       `json:"ram"`
	Redirects []VMRedirect `json:"redirects"`
	Serial    struct {
		ID   string `json:"id"`
		Data string `json:"data"`
		More bool   `json:"more"`
	} `json:"serial"`
	Source struct {
		Checksum   string   `json:"checksum"`
		Filesystem []string `json:"filesystem"`
		Icon       string   `json:"icon"`
		Job        string   `json:"job"`
		Name       string   `json:"name"`
		Type       string   `json:"type"`
	} `json:"source"`
	StateLog string `json:"stateLog"`
	Status   string `json:"status"`
	Summary  string `json:"summary"`
	URL      string `json:"url"`
	Version  string `json:"version"`
}

// VMNetwork ..
type VMNetwork struct {
	Gateway string       `json:"gateway"`
	HTTP    []VMRouteMap `json:"http"`
	HTTPS   []VMRouteMap `json:"https"`
	TCP     []VMRouteMap `json:"tcp"`
	UDP     []VMRouteMap `json:"udp"`
	IP      string       `json:"ip"`
	Mask    string       `json:"mask"`
	Name    string       `json:"name"`
}

// VMRedirect ..
type VMRedirect struct {
	Address string `json:"address"`
	Source  string `json:"source"`
}

// VMRouteMap ..
type VMRouteMap struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

// VorteilConfiguration ..
type VorteilConfiguration struct {
	Args      string              `json:"args"`
	Binary    string              `json:"binary"`
	Env       StringMap           `json:"env"`
	Info      ConfigInfo          `json:"info"`
	Networks  []NetworkInterface  `json:"networks"`
	NFS       []ConfigNFS         `json:"nfs"`
	Redirects StringMap           `json:"redirects"`
	System    SystemConfiguration `json:"system"`
	VM        ConfigVM            `json:"vm"`
}
