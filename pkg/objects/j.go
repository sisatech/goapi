package objects

// Job ..
type Job struct {
	Description      string      `json:"description"`
	ID               string      `json:"id"`
	LogFilePath      string      `json:"logFilePath"`
	logPlainFilePath string      `json:"logPlainFilePath"`
	Name             string      `json:"name"`
	Progress         JobProgress `json:"progress"`
}

// JobProgress ..
type JobProgress struct {
	Error    string  `json:"error"`
	Finished int     `json:"finished"`
	Progress float64 `json:"progress"`
	Started  int     `json:"started"`
	Status   string  `json:"status"`
	Total    float64 `json:"total"`
	Units    string  `json:"units"`
}

type JobEdges struct {
	Cursor string `json:"cursor"`
	Node   Job    `json:"node"`
}

// JobsConnection ..
type JobsConnection struct {
	Edges    []JobEdges `json:"edges"`
	PageInfo PageInfo   `json:"pageInfo"`
}
