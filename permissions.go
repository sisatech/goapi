package goapi

// ApplyACLRule adds read,write or exec to the object.
func (c *Client) ApplyACLRule(id, action, group string) error {
	req := c.NewRequest(`mutation ($id: String!, $group: String!, $action: String!){
		applyACLRule(id: $id, group: $group, action: $action){
			id
		}
	}`)

	req.Var("id", id)
	req.Var("group", group)
	req.Var("action", action)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}
