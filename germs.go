package goapi

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/sisatech/goapi/pkg/objects"
	"github.com/sisatech/vcli/pkg/util/file"
)

// InjectArgs ...
type InjectArgs struct {
	ID                 string
	Type               string
	IsRemoteSource     bool
	IsCompressed       bool
	DecompressedLength int
	ModTime            time.Time
	Path               string
	IsDir              bool
	Replace            bool
	Payload            file.File
}

// URL ...
func (c *Client) URL(typeof, uri string) string {
	var path string
	switch typeof {
	case "build":
		path = "/api/build/"
	case "pack":
		path = "/api/pack/"
	case "unpack":
		path = "/api/unpack/"
	case "push":
		path = "/api/push/"
	case "provision":
		path = "/api/provision"
	}

	return c.cfg.Address + path + uri

}

// AnalyzeQuery ...
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

// DownloadGerm ... takes a save path to download the file, signed uri and type you are downloading and calls a get request for the file that it returns.
func (c *Client) DownloadGerm(save, uri, typeof string) error {

	resp, err := c.Get(c.URL(typeof, uri))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(save)
	if err != nil {
		return err
	}

	defer out.Close()

	io.Copy(out, resp.Body)
	return nil
}

// InjectConfiguration ...
func (c *Client) InjectConfiguration(op *objects.GerminateOperation, id, typeof string, f file.File, replace bool) error {
	defer f.Close()
	return c.Inject(op, &InjectArgs{
		ID:      id,
		Type:    "configuration",
		ModTime: f.ModTime(),
		Payload: f,
		Replace: replace,
	}, typeof)
}

// InjectDirectory ...
func (c *Client) InjectDirectory(op *objects.GerminateOperation, id, path, typeof string, f file.File) error {
	defer f.Close()
	return c.Inject(op, &InjectArgs{
		ID:      id,
		Type:    "archive",
		Path:    path,
		ModTime: f.ModTime(),
		Payload: f,
	}, typeof)
}

// InjectFile ...
func (c *Client) InjectFile(op *objects.GerminateOperation, id, path, typeof string, f file.File) error {
	defer f.Close()
	return c.Inject(op, &InjectArgs{
		ID:      id,
		Path:    path,
		Type:    "file",
		ModTime: f.ModTime(),
		Payload: f,
	}, typeof)
}

// InjectPackage ...
func (c *Client) InjectPackage(op *objects.GerminateOperation, id, typeof string, f file.File) error {
	defer f.Close()
	return c.Inject(op, &InjectArgs{
		ID:      id,
		Type:    "package",
		Payload: f,
	}, typeof)
}

// InjectIcon ...
func (c *Client) InjectIcon(op *objects.GerminateOperation, id, typeof string, f file.File) error {
	defer f.Close()
	return c.Inject(op, &InjectArgs{
		ID:      id,
		Type:    "icon",
		ModTime: f.ModTime(),
		Payload: f,
	}, typeof)
}

// InjectBinary ...
func (c *Client) InjectBinary(op *objects.GerminateOperation, id, typeof string, f file.File) error {
	defer f.Close()
	return c.Inject(op, &InjectArgs{
		ID:      id,
		Type:    "binary",
		ModTime: f.ModTime(),
		Payload: f,
	}, typeof)
}

// InjectDir ...
func (c *Client) InjectDir(op *objects.GerminateOperation, id, path, typeof string) error {
	return c.Inject(op, &InjectArgs{
		ID:    id,
		Type:  "file",
		Path:  path,
		IsDir: true,
	}, typeof)
}

// Inject parses a germinate operation and inject arguments and sends Post request of the file with configuration details on what it is replacing.
func (c *Client) Inject(op *objects.GerminateOperation, args *InjectArgs, typeof string) error {
	defer func() {
		if args.Payload != nil {
			args.Payload.Close()
		}
	}()

	r, err := http.NewRequest(http.MethodPost, c.URL(typeof, op.URI), args.Payload)
	if err != nil {
		return err
	}

	if args.Payload != nil {
		r.Header.Set("Content-Length", fmt.Sprintf("%v", args.Payload.Size()))
		r.ContentLength = int64(args.Payload.Size())
	}

	r.Header.Set("Injection-ID", args.ID)
	r.Header.Set("Injection-Type", args.Type)
	r.Header.Set("Injection-IsRemoteSource", fmt.Sprintf("%v", args.IsRemoteSource))
	r.Header.Set("Injection-IsCompressed", fmt.Sprintf("%v", args.IsCompressed))
	r.Header.Set("Injection-IsDir", fmt.Sprintf("%v", args.IsDir))
	r.Header.Set("Injection-Replace", fmt.Sprintf("%v", args.Replace))

	if args.Path != "" {
		r.Header.Set("Injection-Path", args.Path)
	}

	r.Header.Set("Injection-DecompressedLength", fmt.Sprintf("%v", args.DecompressedLength))
	r.Header.Set("Injection-ModTime", args.ModTime.Format(time.RFC1123))

	resp, err := c.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("status [%v] %s", resp.StatusCode, data)
	}

	return nil
}

// Build ... takes a germ string to build from, a disk format, a kerneltype and injections. Returns a GerminateOperation which has a Job Object and an URI. The URI is used to post your injections to {vorteil daemon uri}/api/build/{uri}.
// If you added injections the graphql query will not complete until the injections are added.
func (c *Client) Build(germ, format, kernel string, injections []string) (*objects.GerminateOperation, error) {

	if format == "" {
		format = "vmdk"
	}

	req := c.NewRequest(`mutation($germ: GermString!, $format: String!, $kernel: String!, $injections: [String]){
		build(germ: $germ, diskFormat: $format, kernelType: $kernel, injections: $injections){
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
	req.Var("format", format)
	req.Var("kernel", kernel)
	req.Var("injections", injections)

	type responseContainer struct {
		GerminateOperation objects.GerminateOperation `json:"build"`
	}

	buildWrapper := new(responseContainer)

	if err := c.client.Run(c.ctx, req, &buildWrapper); err != nil {
		return nil, err
	}

	return &buildWrapper.GerminateOperation, nil

}
