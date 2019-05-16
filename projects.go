package goapi

// ImportSharedObjects gathers libraries from linux to make the binary work correctly. Returns an array of libraries imported or an error if unsuccessfull.
func (c *Client) ImportSharedObjects(project, variant string) ([]string, error) {

	if variant == "" {
		variant = "default"
	}

	req := c.NewRequest(`mutation($variant: String!, $project: String!){
		importSharedObjects(project: $project, variant: $variant)
	}`)

	req.Var("project", project)
	req.Var("variant", variant)

	type responseContainer struct {
		ImportSharedObjects []string `json:"importSharedObjects"`
	}

	sharedObjsWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &sharedObjsWrapper); err != nil {
		return nil, err
	}

	return sharedObjsWrapper.ImportSharedObjects, nil
}

// AddProject function adds a project to the vorteil list. Returns an error if it fails.
func (c *Client) AddProject(path string) error {
	req := c.NewRequest(`query($path: String!){
		addProject{
			path
		}
	}`)

	req.Var("path", path)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}
