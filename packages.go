package goapi

import (
	"fmt"

	"github.com/sisatech/goapi/pkg/objects"
)

// PackageInfo ...
func (c *Client) PackageInfo(path string) (*objects.PackageInfo, error) {

	req := c.NewRequest(`query($path: String!){
		packageInfo(path: $path){
			components {
				binary {
				  modtime
				  name
				  size
				}
				filesystem {
				  modtime
				  name
				  size
				}
				vcfg {
				  modtime
				  name
				  size
				}
			  }
			  configurationDetails {
				app
				author
				binaryArgs
				cpus
				description
				diskSize
				kernel
				memory
				summary
				totalNICs
				url
				version
			  }
			  files
			  id
			  timestamp
		}
	}`)

	req.Var("path", path)

	type responseContainer struct {
		PackageInfo objects.PackageInfo `json:"packageInfo"`
	}

	piWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &piWrapper); err != nil {
		return nil, err
	}

	return &piWrapper.PackageInfo, nil
}

// PackageConfig ...
func (c *Client) PackageConfig(bucket, app, ref string) (*objects.PackageConfig, error) {

	req := c.NewRequest(`query($bucket: String!, $app: String!, $ref: String!){
		packageConfig(bucket: $bucket, app: $app, ref: $ref){
			info{
				app
				author
				binaryArgs
				cpus
				description
				diskSize
				kernel
				memory
				summary
				totalNICs
				url
				version
			}
			raw
		}
	}`)

	req.Var("bucket", bucket)
	req.Var("app", app)
	req.Var("ref", ref)

	type responseContainer struct {
		PackageConfig objects.PackageConfig `json:"packageConfig"`
	}

	pcWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &pcWrapper); err != nil {
		return nil, err
	}

	return &pcWrapper.PackageConfig, nil
}

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

// PackageConfig ..
func (c *Client) PackageConfig(bucket, app, ref string) (*objects.PackageConfig, error) {

	req := c.NewRequest(fmt.Sprintf(`
		query {
			packageConfig(bucket: "%s", app: "%s", ref: "%s") {
				raw
				info {
					app
					author
					binaryArgs
					cpus
					description
					diskSize
					kernel
					memory
					summary
					totalNICs
					url
					version
				}
			}
		}
	`, bucket, app, ref))

	type responseContainer struct {
		PackageConfig objects.PackageConfig `json:"packageConfig"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.PackageConfig, nil
}

// PackageInfo ..
func (c *Client) PackageInfo(path string) (*objects.PackageInfo, error) {

	req := c.NewRequest(fmt.Sprintf(`
		query {
			packageInfo(path: "%s") {
				id
				files
				timestamp
				components {
					vcfg {
						name
						size
						modtime
					}
					binary {
						name
						size
						modtime
					}
					filesystem {
						name
						size
						modtime
					}
				}
				configurationDetails {
					app
					url
					cpus
					author
					kernel
					memory
					summary
					version
					diskSize
					totalNICs
					binaryArgs
					description
				}
			}
		}
	`, path))

	type responseContainer struct {
		PackageInfo objects.PackageInfo `json:"packageInfo"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.PackageInfo, nil
}
