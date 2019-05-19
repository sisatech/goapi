package goapi

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
