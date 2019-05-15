package objects

// Package ..
type Package struct {
	File              FragmentReadExposedOnly `json:"file"`
	Icon              FragmentReadExposedOnly `json:"icon"`
	ID                string                  `json:"id"`
	Tag               string                  `json:"tag"`
	UploadedTimeplate int                     `json:"uploadedTimeplate"`
}

// PackagesConnection ..
type PackagesConnection struct {
	Edges struct {
		Cursor string    `json:"cursor"`
		Node   []Package `json:"node"`
	} `json:"edges"`
	PageInfo PageInfo `json:"pageInfo"`
}

// PageInfo ..
type PageInfo struct {
	EndCursor       string `json:"endCursor"`
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor"`
}

// Platform ..
type Platform struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
