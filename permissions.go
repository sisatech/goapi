package goapi

import (
	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

// ListSingletons ...
func (c *Client) ListSingletons() ([]string, error) {
	req := c.NewRequest(`query{
		listSingletons
	}`)

	type responseContainer struct {
		ListSingletons []string `json:"listSingletons"`
	}

	singletonsWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &singletonsWrapper); err != nil {
		return nil, err
	}

	return singletonsWrapper.ListSingletons, nil
}

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

// DeleteACLRule removes read,write or exec to the object
func (c *Client) DeleteACLRule(id, action, group string) error {
	req := c.NewRequest(`mutation($id: String!, $group: String!, $action: String!){
		deleteACLRule(id: $id, group: $group, action: $action){
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

// ListACLsQuery ..
func (c *Client) ListACLsQuery(id string) ([]*objects.ACL, error) {

	req := graphql.NewRequest(`
                query($id: String!) {
                        listACLRules(id: $id) {
                                acls {
                                        group
                                        action
                                }
                        }
                }
        `)
	req.Var("id", id)

	type responseContainer struct {
		ListACLRules struct {
			ACLs []*objects.ACL `json:"acls"`
		} `json:"listACLRules"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.ListACLRules.ACLs, nil
}

// GetSingletonID ..
func (c *Client) GetSingletonID(name string) (string, error) {

	req := graphql.NewRequest(`
		query($type: Singletons!) {
			getSingletonID(type:$type)
		}
	`)
	req.Var("type", name)

	type responseContainer struct {
		GetSingletonID string `json:"getSingletonID"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return "", err
	}

	return resp.GetSingletonID, nil
}
