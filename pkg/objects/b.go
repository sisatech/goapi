package objects

// Bucket ..
type Bucket struct {
	App           App            `json:"app"`
	AppsList      AppsConnection `json:"appsList"`
	Authorization Authorization  `json:"authorization"`
	Icon          Fragment       `json:"icon"`
	Name          string         `json:"name"`
}
