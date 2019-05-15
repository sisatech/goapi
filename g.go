package goapi

import (
	"fmt"

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

	type responseHolder struct {
		GermConfig objects.VorteilConfiguration `json:"germConfig"`
	}

	response := new(responseHolder)

	err := c.GraphQL().Run(c.Context(), req, &response)
	if err != nil {
		return nil, err
	}

	return &response.GermConfig, nil
}
