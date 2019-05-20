package goapi

import (
	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

// ServerConfig ....
func (c *Client) ServerConfig() (*objects.VorteilConfiguration, error) {
	req := c.NewRequest(`query{
		serverConfig{
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
	}`)

	type responseContainer struct {
		ServerConfig objects.VorteilConfiguration `json:"serverConfig"`
	}

	scWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &scWrapper); err != nil {
		return nil, err
	}

	return &scWrapper.ServerConfig, nil
}

// LogFile ...
func (c *Client) LogFile() (string, error) {
	req := c.NewRequest(`query{
		logFile
	}`)

	type responseContainer struct {
		LogFile string `json:"logFile"`
	}

	logWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &logWrapper); err != nil {
		return "", err
	}

	return logWrapper.LogFile, nil
}

// ListKernels ..
func (c *Client) ListKernels() ([]objects.KernelVersion, error) {
	req := c.NewRequest(`query{
		listKernels{
			release
			source
			type
			version
		}
	}`)

	type responseContainer struct {
		ListKernels []objects.KernelVersion `json:"listKernels"`
	}

	kernelsWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &kernelsWrapper); err != nil {
		return nil, err
	}

	return kernelsWrapper.ListKernels, nil
}

// ListDefaultsOptions ..
func (c *Client) ListDefaultsOptions() (*objects.Lists, error) {
	req := c.NewRequest(`query{
		listDefaultsOptions{
			kernels{
				release
				source
				type
				version
			}
			platforms
		}
	}`)

	type responseContainer struct {
		ListDefaultsOptions objects.Lists `json:"listDefaultsOptions"`
	}

	defaultsWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &defaultsWrapper); err != nil {
		return nil, err
	}

	return &defaultsWrapper.ListDefaultsOptions, nil
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
