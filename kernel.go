package goapi

import (
	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

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

// UpdateKernels ...
func (c *Client) UpdateKernels() error {
	req := c.NewRequest(`mutation{
		updateKernels
	}`)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// SetDefault ...
func (c *Client) SetDefault(kernel, platform string) error {
	req := c.NewRequest(`mutation($kernel: String, $platform: String){
		setDefault(kernel: $kernel, platform: $platform){
			kernel
			platform
		}
	}`)

	if kernel != "" {
		req.Var("kernel", kernel)
	}

	if platform != "" {
		req.Var("platform", platform)
	}

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil

}

// RemoveKernel ...
func (c *Client) RemoveKernel(version string) error {
	req := c.NewRequest(`mutation($version: String!){
		removeKernel(version: $version)
	}`)

	req.Var("version", version)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// DownloadKernel ...
func (c *Client) DownloadKernel(version string) error {
	req := c.NewRequest(`mutation($version: String!){
		downloadKernel(version: $version){
			release
			source
			type
			version
		}
	}`)

	req.Var("version", version)

	if err := c.client.Run(c.ctx, req, nil); err != nil {
		return err
	}

	return nil
}
