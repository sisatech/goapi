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

// MachinesManager provides access to 'Machine' APIs
type MachinesManager struct {
	c           *Client
	environment *Repository
}

// WithEnvironment ..
func (m *MachinesManager) WithEnvironment(r *Repository) *MachinesManager {
	out := new(MachinesManager)
	*out = *m
	out.environment = r
	return out
}

// Provision a virtual machine.
func (m *MachinesManager) Provision(args *ProvisionArguments) (*ProvisionOperation, error) {

	if args.Injections == nil {
		args.Injections = make([]string, 0)
	}
	for _, a := range args.Injections {
		a = fmt.Sprintf(`"%s"`, a)
	}

	argsStr := make([]string, 0)
	argsStr = append(argsStr, fmt.Sprintf(`germ: "%s"`, args.Germ))
	if args.InstanceName != "" {
		argsStr = append(argsStr, fmt.Sprintf(`name: "%s"`, args.InstanceName))
	}
	if args.KernelType != "" {
		argsStr = append(argsStr, fmt.Sprintf(`kernelType: "%s"`, args.KernelType))
	}
	if args.Platform != "" {
		argsStr = append(argsStr, fmt.Sprintf(`platform: "%s"`, args.Platform))
	}
	argsStr = append(argsStr, fmt.Sprintf(`start: %v`, args.PoweredOn))
	if len(args.Injections) != 0 {
		argsStr = append(argsStr, fmt.Sprintf("injections: [%s]", strings.Join(args.Injections, ", ")))
	}

	req := m.environment.newRequest(fmt.Sprintf(`
		mutation {
			provision (%s) {
				uri
				job {
					id
				}
			}
		}
	`, strings.Join(argsStr, ", ")))

	type responseContainer struct {
		Provision objects.CompoundProvisionResponse `json:"provision"`
	}

	resp := new(responseContainer)
	err := m.c.graphql.Run(m.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &ProvisionOperation{
		c:     m.c,
		host:  m.environment.host,
		jobID: resp.Provision.Job.ID,
		uri:   resp.Provision.URI,
	}, nil
}

// Inject ..
func (p *ProvisionOperation) Inject(key string, itype InjectionType, value io.Reader, headers http.Header) error {

	url := fmt.Sprintf("%s/api/provision/%s", p.host, p.uri)
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

// List all virtual machines.
func (m *MachinesManager) ListMachines(cursor *Cursor) (*VirtualMachineList, error) {

	var vd, v string
	if cursor != nil {
		vd, v = cursor.Strings()
	}

	req := m.environment.newRequest(fmt.Sprintf(`
		query%s {
			listMachines%s {
				edges {
					node {
						id
						name
					}
					cursor
				}
				pageInfo {
					endCursor
					startCursor
					hasNextPage
					hasPreviousPage
				}
			}
		}
	`, vd, v))

	if cursor != nil {
		cursor.AddToRequest(req)
	}

	type responseContainer struct {
		ListMachines objects.VMsConnection `json:"listMachines"`
	}
	resp := new(responseContainer)
	err := m.c.graphql.Run(m.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	out := &VirtualMachineList{
		PageInfo: resp.ListMachines.PageInfo,
		Items:    make([]VirtualMachineListItem, 0),
	}
	for _, v := range resp.ListMachines.Edges {
		out.Items = append(out.Items, VirtualMachineListItem{
			Cursor: v.Cursor,
			VirtualMachine: VirtualMachine{
				id:   v.Node.ID,
				name: v.Node.Name,
			},
		})
	}

	return out, nil
}

// Get a virtual machine.
func (m *MachinesManager) Get(id string) (*VirtualMachine, error) {

	req := m.environment.newRequest(fmt.Sprintf(`
		query {
			vm(id: "%s") {
				id
				name
			}
		}
	`, id))

	type responseContainer struct {
		VM objects.VM `json:"vm"`
	}

	resp := new(responseContainer)
	err := m.c.graphql.Run(m.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &VirtualMachine{
		mgr:  m,
		id:   resp.VM.ID,
		name: resp.VM.Name,
	}, nil
}

// WaitUntilFinished ..
func (p *ProvisionOperation) WaitUntilFinished() error {

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
		err := p.c.graphql.Run(context.TODO(), req, &resp)
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

// ProvisionOperation ..
type ProvisionOperation struct {
	c     *Client
	jobID string
	uri   string
	host  string
}

// ProvisionArguments ..
type ProvisionArguments struct {
	Germ         string
	InstanceName string
	PoweredOn    bool
	Platform     string
	KernelType   KernelType
	Injections   []string
}
