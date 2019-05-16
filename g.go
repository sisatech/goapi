package goapi

import (
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

// GermConfigQuery ..
func (c *Client) GermConfigQuery(germ string) (*objects.VorteilConfiguration, error) {

	req := c.NewRequest(fmt.Sprintf(`
                query ($germ: String!) {
                        germConfig (germ:$germ) {
                                args
                                binary
                                env {
                                        tuples {
                                                key
                                                value
                                        }
                                }
                                info {
                                        author
                                        description
                                        name
                                        summary
                                        url
                                        version
                                }
                                networks {
                                        disableTCPSegmentationOffload
                                        gateway
                                        http
                                        https
                                        ip
                                        mask
                                        mtu
                                        tcp
                                        udp
                                }
                                nfs {
                                        mountPoint
                                        server
                                }
                                redirects {
                                        tuples {
                                                key
                                                value
                                        }
                                }
                                system {
                                        delay
                                        diskCache
                                        dns
                                        hostname
                                        maxFDs
                                        outputFormat
                                        pages4k
                                        stdoutMode
                                }
                                vm {
                                        cpus
                                        diskSize
                                        inodes
                                        kernel
                                        ram
                                }
                        }
                }
        `))

	req.Var("germ", germ)

	type responseContainer struct {
		GermConfig objects.VorteilConfiguration `json:"germConfig"`
	}

	response := new(responseContainer)

	err := c.client.Run(c.Context(), req, &response)
	if err != nil {
		return nil, err
	}

	return &response.GermConfig, nil
}

// GetDefaultQuery ..
func (c *Client) GetDefaultQuery() (*objects.Defaults, error) {

	req := graphql.NewRequest(`
	query {
		getDefault {
			kernel
			platform
		}
	}
	`)

	type responseContainer struct {
		GetDefault objects.Defaults `json:"getDefault"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.GetDefault, nil
}
