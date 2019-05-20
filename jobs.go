package goapi

import (
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

// JobQuery ..
func (c *Client) JobQuery(id string) (*objects.Job, error) {

	req := graphql.NewRequest(`
                query($id: String!) {
                        job(id:$id) {
                                id
                                name
                                progress {
                                        error
                                        total
                                        units
                                        status
                                        started
                                        finished
                                        progress
                                }
                                description
                                logFilePath
                                logPlainFilePath
                        }
                }
        `)
	req.Var("id", id)

	type responseContainer struct {
		Job objects.Job `json:"job"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Job, nil
}

// JobsQuery ..
func (c *Client) JobsQuery(cursor *CursorArgs) (*objects.JobsConnection, error) {

	variableDeclarations, variables := cursor.Strings()
	req := graphql.NewRequest(fmt.Sprintf(`
                query%s {
                        jobs%s {
                                edges {
                                        node {
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
                                        cursor
                                }
                                pageInfo {
                                        endCursor
                                        hasNextPage
                                        startCursor
                                        hasPreviousPage
                                }
                        }
                }
        `, variableDeclarations, variables))
	cursor.AddToRequest(req)

	type responseContainer struct {
		Jobs objects.JobsConnection `json:"jobs"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Jobs, nil
}

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
