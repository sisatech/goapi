package objects

// Linter ..
type Linter struct {
	Line int    `json:"line"`
	Lvl  string `json:"lvl"`
	Msg  string `json:"msg"`
}

// Lists ..
type Lists struct {
	Kernels []struct {
		Release string `json:"release"`
		Source  string `json:"source"`
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"kernels"`
	Platforms []string `json:"platforms"`
}
