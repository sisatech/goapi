package goapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

// Push ..
func (r *Repository) Push(args *PushArguments) (*PushOperation, error) {

	if args.Injections == nil {
		args.Injections = make([]string, 0)
	}
	for _, a := range args.Injections {
		a = fmt.Sprintf(`"%s"`, a)
	}

	argStrings := make([]string, 0)
	argStrings = append(argStrings, fmt.Sprintf(`germ: "%s"`, args.Germ))
	argStrings = append(argStrings, fmt.Sprintf(`bucket: "%s"`, args.DestinationBucket))
	argStrings = append(argStrings, fmt.Sprintf(`app: "%s"`, args.DestinationApp))

	if args.RepositoryName != "" {
		argStrings = append(argStrings, fmt.Sprintf(`node: "%s"`, args.RepositoryName))
	}

	argStrings = append(argStrings, fmt.Sprintf("injections: [%s]", strings.Join(args.Injections, ", ")))
	pushArgs := strings.Join(argStrings, ", ")

	req := r.newRequest(fmt.Sprintf(`
		mutation {
			push(%s) {
				uri
				job {
					id
				}
			}
		}
	`, pushArgs))

	type responseContainer struct {
		Push objects.GerminateOperation `json:"push"`
	}
	resp := new(responseContainer)
	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	out := new(PushOperation)
	out.graphql = r.mgr.c.graphql
	out.jobID = resp.Push.Job.ID
	out.uri = resp.Push.URI

	if r.name != "local" {
		out.host = r.host
	} else {
		out.host = fmt.Sprintf("%s%s", r.mgr.c.protocol, r.mgr.c.cfg.Address)
	}

	return out, nil
}

// PushArguments ..
type PushArguments struct {
	Germ              string
	DestinationBucket string
	DestinationApp    string
	RepositoryName    string
	Injections        []string
}

// PushOperation ..
type PushOperation struct {
	graphql *graphql.Client
	jobID   string
	uri     string
	host    string
}

// WaitUntilFinished ..
func (p *PushOperation) WaitUntilFinished() error {

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
        `, p.jobID))

	type responseContainer struct {
		Job objects.Job `json:"job"`
	}
	resp := new(responseContainer)

	for {
		err := p.graphql.Run(context.TODO(), req, &resp)
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

// InjectionType ..
type InjectionType string

const (
	BinaryInjection        = "binary"
	ConfigurationInjection = "configuration"
	IconInjection          = "icon"
	FileInjection          = "file"
	ArchiveInjection       = "archive"
	PackageInjection       = "package"
)

// Inject ..
func (p *PushOperation) Inject(key string, itype InjectionType, value io.Reader, headers http.Header) error {

	url := fmt.Sprintf("%s/api/push/%s", p.host, p.uri)
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
