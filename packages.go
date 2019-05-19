package goapi

import (
	"github.com/sisatech/goapi/pkg/objects"
)

// Unpack ...
func (c *Client) Unpack(germ string, injections []string) (*objects.GerminateOperation, error) {

	req := c.NewRequest(`mutation($germ: GermString!, $injections: [String]){
		unpack(germ: $germ, injections: $injections){
			job{
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
			uri
		}
	}`)

	req.Var("germ", germ)
	req.Var("injections", injections)

	type responseContainer struct {
		Unpack objects.GerminateOperation `json:"unpack"`
	}

	unpackWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &unpackWrapper); err != nil {
		return nil, err
	}

	return &unpackWrapper.Unpack, nil
}

// Pack ...
func (c *Client) Pack(germ string, compressionLevel int, injections []string) (*objects.GerminateOperation, error) {

	req := c.NewRequest(`mutation($germ: GermString!, $compression: Int, $injections: [String]){
		pack(germ: $germ, compressionLevel: $compression, injections: $injections){
			job {
				id
				description
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
			uri
		}
	}`)

	req.Var("germ", germ)
	req.Var("compression", compressionLevel)
	req.Var("injections", injections)

	type responseContainer struct {
		GerminateOperation objects.GerminateOperation `json:"pack"`
	}

	packWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &packWrapper); err != nil {
		return nil, err
	}

	return &packWrapper.GerminateOperation, nil
}
