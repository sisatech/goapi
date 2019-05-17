package goapi

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
