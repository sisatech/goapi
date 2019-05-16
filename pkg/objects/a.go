package objects

// Authorization ..
type ACL struct {
	Action string `json:"action"`
	Group  string `json:"group"`
}

// App ..
type App struct {
	Authorization Authorization      `json:"authorization"`
	Latest        Package            `json:"latest"`
	Name          string             `json:"name"`
	Version       Package            `json:"version"`
	VersionsList  PackagesConnection `json:"versionsList"`
}

// AppEdges ..
type AppEdges struct {
	Cursor string `json:"cursor"`
	Node   App    `json:"node"`
}

// AppsConnection ..
type AppsConnection struct {
	Edges    []AppEdges `json:"edges"`
	PageInfo PageInfo   `json:"pageInfo"`
}

// Authorization ..
type Authorization struct {
	ACLS  []ACL  `json:"acls"`
	ID    string `json:"id"`
	Owner string `json:"owner"`
}
