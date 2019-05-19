package objects

// Fragment ...
type Fragment struct {
	DownloadURL string `json:"downloadURL"`
	ID          string `json:"id"`
	MD5         string `json:"md5"`
	UploadURL   string `json:"uploadURL"`
}

// FragmentReadExposedOnly ..
type FragmentReadExposedOnly Fragment

// FileInfo ..
type FileInfo struct {
	ModTime int    `json:"modTime"`
	Name    string `json:"name"`
	Size    int    `json:"size"`
}
