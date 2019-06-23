package objects

// Defaults ..
type Defaults struct {
	Kernel   string `json:"kernel"`
	Platform string `json:"platform"`
}

// DetailPlatform ..
type DetailPlatform struct {
	Debug bool                   `json:"dbug"`
	More  map[string]interface{} `json:"more"`
}

// DiskAnalysis ..
type DiskAnalysis struct {
	Filesystem DiskFilesystem `json:"fileSystem"`
}

// DiskFilesystem ..
type DiskFilesystem struct {
	Contents []FSInfo `json:"contents"`
}

// FSInfo ..
type FSInfo struct {
	AccessTime string `json:"accessTime"`
	IsDir      bool   `json:"isDir"`
	ModTime    string `json:"modTime"`
	Mode       int    `json:"mode"`
	Path       string `json:"path"`
	Size       int    `json:"size"`
}
