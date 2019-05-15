package objects

// StringMap ..
type StringMap struct {
	Tuples []MapTuple `json:"tuples"`
}

// SystemConfiguration ..
type SystemConfiguration struct {
	Delay        string   `json:"delay"`
	DiskCache    string   `json:"diskCache"`
	DNS          []string `json:"dns"`
	Hostname     string   `json:"hostname"`
	MaxFDs       int      `json:"maxFDs"`
	OutputFormat string   `json:"outputFormat"`
	Pages4K      bool     `json:"pages4k"`
	StdoutMode   string   `json:"stdoutMode"`
}
