package objects

// CompoundProvisionResponse ..
type CompoundProvisionResponse struct {
	ID  string `json:"id"`
	Job Job    `json:"job"`
	URI string `json:"uri"`
}

// ConfigInfo ..
type ConfigInfo struct {
	Author      string `json:"author"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	URL         string `json:"url"`
	Version     string `json:"version"`
}

// ConfigNFS ..
type ConfigNFS struct {
	MountPoint string `json:"mountPoint"`
	Server     string `json:"server"`
}

// ConfigVM ..
type ConfigVM struct {
	CPUs     int    `json:"cpus"`
	DiskSize string `json:"diskSize"`
	INodes   int    `json:"inodes"`
	Kernel   string `json:"kernel"`
	RAM      string `json:"ram"`
}
