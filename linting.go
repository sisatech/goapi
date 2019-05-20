package goapi

import (
	"github.com/sisatech/goapi/pkg/objects"
)

// ProjectLinter ..
func (c *Client) ProjectLinter(path string) ([]objects.Linter, error) {
	req := c.NewRequest(`query($path: String!){
		projectLinter(path: $path){
			line
			lvl
			msg
		}
	}`)

	req.Var("path", path)

	type responseContainer struct {
		ProjectLinter []objects.Linter `json:"projectLinter"`
	}

	pLinter := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &pLinter); err != nil {
		return nil, err
	}

	return pLinter.ProjectLinter, nil
}

// VCFGLinter
func (c *Client) VCFGLinter(path string) ([]objects.Linter, error) {
	req := c.NewRequest(`query($path: String!){
		vcfgLinter(path: $path){
			line
			lvl
			msg
		}	
	}`)

	req.Var("path", path)

	type responseContainer struct {
		VCFGLinter []objects.Linter `json:"vcfgLinter"`
	}

	vcfgLinter := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &vcfgLinter); err != nil {
		return nil, err
	}

	return vcfgLinter.VCFGLinter, nil
}
