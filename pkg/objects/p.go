package objects

// Package ..
type Package struct {
	File              PackageFragment `json:"file"`
	Icon              PackageFragment `json:"icon"`
	ID                string          `json:"id"`
	Tag               string          `json:"tag"`
	UploadedTimeplate int             `json:"uploadedTimeplate"`
}

// PackageConfig ..
type PackageConfig struct {
	Info struct {
		App         string `json:"app"`
		Author      string `json:"author"`
		BinaryArgs  string `json:"binaryArgs"`
		CPUs        int    `json:"cpus"`
		Description string `json:"description"`
		DiskSize    int    `json:"diskSize"`
		Kernel      string `json:"kernel"`
		Memory      int    `json:"memory"`
		Summary     string `json:"summary"`
		TotalNICs   int    `json:"totalNICs"`
		URL         string `json:"url"`
		Version     string `json:"version"`
	} `json:"info"`
	Raw string `json:"raw"`
}

// PackageComponents ...
type PackageComponents struct {
	Binary     FileInfo `json:"binary"`
	FileSystem FileInfo `json:"filesystem"`
	VCFG       FileInfo `json:"vcfg"`
}

// PackageInfo ...
type PackageInfo struct {
	Components           PackageComponents `json:"components"`
	ConfigurationDetails PackageConfig     `json:"configurationDetails"`
	Files                []string          `json:"files"`
	ID                   string            `json:"id"`
	Timestamp            int               `json:"timestamp"`
}

// PackagesEdge ..
type PackagesEdge struct {
	Cursor string  `json:"cursor"`
	Node   Package `json:"node"`
}

// PackagesConnection ..
type PackagesConnection struct {
	Edges    []PackagesEdge `json:"edges"`
	PageInfo PageInfo       `json:"pageInfo"`
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
