package goapi

// DeleteJob ... deletes a job ID from a repository
func (c *Client) DeleteJob(id string) error {
	req := c.NewRequest(`mutation($id: String!){
		deleteJob(id: $id){
			description
			id
			logFilePath
			logPlainFilePath
			name
			progress {
			error
			finished
			progress
			started
			status
			total
			units
			}
		}
  }`)

	req.Var("id", id)
	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// CancelJob ...
func (c *Client) CancelJob(id string) error {

	req := c.NewRequest(`mutation($id: String!){
		cancelJob(id: $id){
			description
			id
			logFilePath
			logPlainFilePath
			name
			progress {
			  error
			  finished
			  progress
			  started
			  status
			  total
			  units
			}
		  }
		}
	}`)

	req.Var("id", id)
	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}
