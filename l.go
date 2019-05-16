package goapi

import (
	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

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
