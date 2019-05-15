package objects

// KernelVersion ..
type KernelVersion struct {
	Release int64  `json:"release"`
	Source  string `json:"source"`
	Type    string `json:"type"`
	Version string `json:"version"`
}
