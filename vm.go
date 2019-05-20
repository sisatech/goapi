package goapi

import (
	"encoding/json"
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/graphqlws"
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
func (c *Client) Provision(germ, platform, kernelType, name string, injects []string, start bool) (objects.CompoundProvisionResponse, error) {
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
		return objects.CompoundProvisionResponse{}, err
	}

	return provisionWrapper.Provision, nil
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

// VMList ..
type VMList struct {
	PageInfo objects.PageInfo
	Items    []VMListItem
}

// VMListItem ..
type VMListItem struct {
	Cursor string
	VM     objects.VM
}

// ListVMs ..
func (c *Client) ListVMs(curs *Cursor) (*VMList, error) {

	var vd, v string
	if curs != nil {
		vd, v = curs.Strings()
	}

	req := graphql.NewRequest(fmt.Sprintf(`
		query%s {
			listVMs%s {
				edges {
					cursor
					node {
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
				}
				pageInfo {
					endCursor
					startCursor
					hasNextPage
					hasPreviousPage
				}
			}
		}
	`, vd, v))
	if curs != nil {
		curs.AddToRequest(req)
	}

	type responseContainer struct {
		ListVMs *objects.VMsConnection `json:"listVMs"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	out := new(VMList)
	out.PageInfo = resp.ListVMs.PageInfo
	out.Items = make([]VMListItem, 0)

	for _, v := range resp.ListVMs.Edges {
		out.Items = append(out.Items, VMListItem{
			Cursor: v.Cursor,
			VM:     v.Node,
		})
	}

	return out, nil
}

// ListVMsSubscription ..
func (c *Client) ListVMsSubscription(dataCallback func(list *VMList,
	errs []graphqlws.GQLError), errCallback func(error)) (*graphqlws.Subscription, error) {

	dc := func(payload *graphqlws.GQLDataPayload) {

		if payload.Data == nil {
			dataCallback(nil, payload.Errors)
			return
		}

		type responseContainer struct {
			Data struct {
				ListVMs objects.VMsConnection `json:"listVMs"`
			} `json:"data"`
		}

		resp := new(responseContainer)
		b, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, resp)
		if err != nil {
			panic(err)
		}

		out := new(VMList)
		out.PageInfo = resp.Data.ListVMs.PageInfo
		out.Items = make([]VMListItem, 0)

		for _, v := range resp.Data.ListVMs.Edges {
			out.Items = append(out.Items, VMListItem{
				Cursor: v.Cursor,
				VM:     v.Node,
			})
		}

		dataCallback(out, payload.Errors)
	}

	subscription, err := c.subscriptions.Subscription(&graphqlws.SubscriptionConfig{
		Query: fmt.Sprintf(`
			subscription {
				listVMs {
					pageInfo {
						endCursor
						startCursor
						hasNextPage
						hasPreviousPage
					}
					edges {
						cursor
						node {
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
					}
				}
			}
		`),
		DataCallback:  dc,
		ErrorCallback: errCallback,
	})
	if err != nil {
		return nil, err
	}

	return subscription, nil
}
