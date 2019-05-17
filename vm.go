package goapi

// PauseVM ...
func (c *Client) PauseVM(id string) error {
	req := c.NewRequest(`mutation($id: String!){
		pauseVM(id: $id){
			id
		}
	}`)

	req.Var("id", id)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// DetachVM takes the manager away from the instance for it to be able to stay running permanently.
func (c *Client) DetachVM(id string) error {
	req := c.NewRequest(`mutation($id: String!){
		detachVM(id: $id){
			id
		}
	}`)

	req.Var("id", id)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil

}

// DeleteVM ...
func (c *Client) DeleteVM(id string) error {
	req := c.NewRequest(`mutation($id: String!){
		deleteVM(id: $id){
			id
		}
	}`)

	req.Var("id", id)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// CreateTemplateFromVM ... from a provisioned vm we take its ID and make a template on cloud platforms GCP, AWS
func (c *Client) CreateTemplateFromVM(id string) error {

	req := c.NewRequest(`mutation($id: ID!){
		createTemplate(id: $id) {
			args
			author
			binary
			cpus
			created
			date
			disk
			download
			env
			hostname
			id
			instance
			kernel
			logFile
			name
			networks {
			  gateway
			  http {
				address
				port
			  }
			  https {
				address
				port
			  }
			  ip
			  mask
			  name
			  tcp {
				address
				port
			  }
			  udp {
				address
				port
			  }
			}
			platform
			ram
			redirects {
			  address
			  source
			}
			serial {
			  cursor
			  data
			  more
			}
			source {
			  checksum
			  filesystem
			  icon
			  job
			  name
			  type
			}
			stateLog
			status
			summary
			url
			version
		}
	}`)

	req.Var("id", id)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil

}
