package goapi

import (
	"github.com/sisatech/goapi/pkg/objects"
)

// VM ...
func (c *Client) VM(id string) (*objects.VM, error) {
	req := c.NewRequest(`query($id: ID!){
		vm(id: $id){
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

	type responseContainer struct {
		VM objects.VM `json:"vm"`
	}

	vmWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &vmWrapper); err != nil {
		return nil, err
	}

	return &vmWrapper.VM, nil
}

// StopVM ...
func (c *Client) StopVM(id string) error {
	req := c.NewRequest(`mutation($id: String!){
		stopVM(id: $id){
			id
		}	
	}`)

	req.Var("id", id)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// StartVM ...
func (c *Client) StartVM(id string) error {
	req := c.NewRequest(`mutation($id: ID!){
		startVM(id: $id){
			id
		}
	}`)

	req.Var("id", id)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Provision ...
func (c *Client) Provision(germ, platform, kernelType, name string, injects []string, start bool) (*objects.CompoundProvisionResponse, error) {

	req := c.NewRequest(`mutation($germ: GermString!, $platform: String!, $kernelType: String, $name: String!, $injects: [String], $start: Boolean){
		provision(germ: $germ, platform: $platform, kernelType: $kernelType, name: $name, injections: $injects, start: $start){
			id
			uri
			job {
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

	req.Var("germ", germ)
	req.Var("platform", platform)
	req.Var("kernelType", kernelType)
	req.Var("name", name)
	req.Var("injections", injects)
	req.Var("start", start)

	type responseContainer struct {
		Provision objects.CompoundProvisionResponse `json:"provision"`
	}

	provisionWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &provisionWrapper); err != nil {
		return nil, err
	}

	return &provisionWrapper.Provision, nil
}

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
