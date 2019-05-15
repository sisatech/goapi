package objects

// FragmentReadExposedOnly ..
type Fragment struct {
	DownloadURL string `json:"downloadURL"`
	ID          string `json:"id"`
	MD5         string `json:"md5"`
	UploadURL   string `json:"uploadURL"`
}

// Fragment ..
type FragmentReadExposedOnly Fragment

// FileInfo ..
type FileInfo struct {
	ModTime int    `json:"modTime"`
	Name    string `json:"name"`
	Size    int    `json:"size"`
}
