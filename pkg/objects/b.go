package objects

// Bucket ..
type Bucket struct {
	App           App            `json:"app"`
	AppsList      AppsConnection `json:"appsList"`
	Authorization Authorization  `json:"authorization"`
	Icon          Fragment       `json:"icon"`
	Name          string         `json:"name"`
}

// BucketsConnection ..
type BucketsConnection struct {
	PageInfo PageInfo      `json:"pageInfo"`
	Edges    []BucketsEdge `json:"edges"`
}

// BucketsEdge ..
type BucketsEdge struct {
	Cursor string `json:"cursor"`
	Node   Bucket `json:"node"`
}
