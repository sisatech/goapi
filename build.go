package goapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sisatech/goapi/pkg/objects"
	"github.com/machinebox/graphql"
)

// BuildManager provides access to build APIs.
type BuildManager struct {
	c           *Client
	environment *Repository
}

// WithEnvironment ..
func (m *BuildManager) WithEnvironment(r *Repository) *BuildManager {
	out := new(BuildManager)
	*out = *m
	out.environment = r
	return out
}

// Build a Vorteil disk image.
func (b *BuildManager) Build(args *BuildArguments) (*BuildOperation, error) {

	if args.Injections == nil {
		args.Injections = make([]string, 0)
	}
	for _, a := range args.Injections {
		a = fmt.Sprintf(`"%s"`, a)
	}

	argStrings := make([]string, 0)
	argStrings = append(argStrings, fmt.Sprintf(`germ: "%s"`, args.Germ))
	if args.DiskFormat != "" {
		argStrings = append(argStrings, fmt.Sprintf(`diskFormat: "%s"`, args.DiskFormat))
	}
	if len(args.Injections) != 0 {
		argStrings = append(argStrings, fmt.Sprintf("injections: [%s]", strings.Join(args.Injections, ", ")))
	}

	req := b.environment.newRequest(fmt.Sprintf(`
                mutation {
                        build(%s) {
                                job {
                                        id
                                }
                                uri
                        }
                }
        `, fmt.Sprintf("%s", strings.Join(argStrings, ", "))))
	type responseContainer struct {
		Build objects.GerminateOperation `json:"build"`
	}
	resp := new(responseContainer)
	err := b.environment.mgr.c.graphql.Run(b.environment.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}
	return &BuildOperation{
		graphql: b.environment.mgr.c.graphql,
		host:    b.environment.host,
		jobID:   resp.Build.Job.ID,
		uri:     resp.Build.URI,
	}, nil
}

// Start the build operation by providing an io.Writer to write the disk image to
func (b *BuildOperation) Start(w io.Writer) error {

	url := fmt.Sprintf("%s/api/build/%s", b.host, b.uri)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with a non-200 status code: %v - %s", resp.StatusCode, resp.Status)
	}

	go func() {
		defer resp.Body.Close()
		_, err := io.Copy(w, resp.Body)
		if err != nil {
			// return 0, err
			// TODO properly report error
			fmt.Println(err)
		}
	}()

	return nil
}

// WaitUntilFinished ..
func (b *BuildOperation) WaitUntilFinished() error {

	req := graphql.NewRequest(fmt.Sprintf(`
                query {
                        job(id: "%s") {
                                logFilePath
                                progress {
                                        status
                                        finished
                                        error
                                }
                        }
                }
        `, b.jobID))

	type responseContainer struct {
		Job objects.Job `json:"job"`
	}
	resp := new(responseContainer)

	for {
		err := b.graphql.Run(context.TODO(), req, &resp)
		if err != nil {
			return err
		}

		if resp.Job.Progress.Finished != 0 {
			if resp.Job.Progress.Error != "" {
				return fmt.Errorf("job error: %s", resp.Job.Progress.Error)
			}
			return nil
		}

		time.Sleep(time.Second * 1)
	}

	return nil
}

// Inject ..
func (b *BuildOperation) Inject(key string, itype InjectionType, value io.Reader, headers http.Header) error {

	url := fmt.Sprintf("%s/api/build/%s", b.host, b.uri)
	req, err := http.NewRequest(http.MethodPost, url, value)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	req.Header.Add("Injection-ID", key)
	req.Header.Add("Injection-Type", string(itype))

	// Allow users to pass as many custom headers as desired, if any at all.
	for k, h := range headers {
		if len(h) == 1 {
			req.Header.Set(k, h[0])
		} else {
			for _, x := range h {
				req.Header.Add(k, x)
			}
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status code: %v - %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// BuildArguments contain the fields used in a Build operation.
// Germ         - a 'germ' string is an unambiguous pointer to a valid target for the
//              'build' operation. It could be the path to a package/project, or
//              an app/version within a repository [see (*App).Germ()]
// Injections   - a list of injection IDs, which the user can reference when
//              injecting various components to the build process (such as
//              configuration settings, files, etc.)
// DiskFormat   - the desired format of the resulting disk image.
// KernelType   - "prod" or "debug"
// Repository   - the repository that will perform the build operation (can be
//              left 'nil'; defaults to 'local').
type BuildArguments struct {
	Germ       string
	Injections []string
	DiskFormat DiskFormat
	KernelType KernelType
}

type KernelType string

const (
	ProductionKernel = KernelType("prod")
	DebugKernel      = KernelType("debug")
)

type DiskFormat string

const (
	DiskFormatRAW                 = DiskFormat("raw")
	DiskFormatGCP                 = DiskFormat("gcp")
	DiskFormatVMDK                = DiskFormat("vmdk")
	DiskFormatSparseVMDK          = DiskFormat("sparse-vmdk")
	DiskFormatStreamOptimizedVMDK = DiskFormat("steam-optimized-vmdk")
	DiskFormatOVA                 = DiskFormat("ova")
	DiskFormatVHD                 = DiskFormat("vhd")
	DiskFormatXVA                 = DiskFormat("xva")
)

// BuildOperation ..
type BuildOperation struct {
	graphql *graphql.Client
	jobID   string
	uri     string
	host    string
}

// AnalyzeDisk ..
func (b *BuildManager) AnalyzeDisk(path string) (*FilesystemInfo, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
		query {
			analyze(path: "%s") {
				fileSystem {
					contents {
						mode
						path
						size
						isDir
						modTime
						accessTime
					}
				}
			}
		}
	`, path))

	type responseContainer struct {
		Analyze objects.DiskAnalysis `json:"analyze"`
	}

	resp := new(responseContainer)
	err := b.c.graphql.Run(b.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	out := new(FilesystemInfo)
	out.Contents = make([]FileInfo, 0)
	for _, x := range resp.Analyze.Filesystem.Contents {
		out.Contents = append(out.Contents, FileInfo{
			FSInfo: x,
		})
	}

	return out, nil
}

// FilesystemInfo ..
type FilesystemInfo struct {
	Contents []FileInfo
}

// FileInfo ..
type FileInfo struct {
	objects.FSInfo
}
