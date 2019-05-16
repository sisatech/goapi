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
	Contents struct {
		AccessTime int64  `json:"accessTime"`
		IsDir      bool   `json:"isDir"`
		ModTime    int64  `json:"modTime"`
		Mode       int    `json:"mode"`
		Path       string `json:"path"`
		Size       int    `json:"size"`
	} `json:"contents"`
}
