package goapi

import (
	"fmt"

	"github.com/sisatech/goapi/pkg/objects"
)

// AnalyzeQuery
func (c *Client) AnalyzeQuery(path string) (*objects.DiskAnalysis, error) {

	req := c.NewRequest(fmt.Sprintf(`
                query ($path: String!) {
                        analyze (path:$path) {
				fileSystem {
					contents {
						accessTime
						isDir
						modTime
						mode
						path
						size
					}
				}
			}
                }
        `))
	req.Var("path", path)

	type responseContainer struct {
		Analyze objects.DiskAnalysis `json:"analyze"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Analyze, nil
}
