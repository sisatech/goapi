package goapi

import (
	"fmt"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

// RepositoriesManager provides access to the repositories APIs.
type RepositoriesManager struct {
	c     *Client
	Local *Repository
}

// Connections lists all connected repositories.
func (r *RepositoriesManager) Connections() ([]Repository, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
		query {
			listNodes {
				name
				host
			}
		}
	`))

	type responseContanier struct {
		ListNodes []objects.Node `json:"listNodes"`
	}
	resp := new(responseContanier)
	err := r.c.graphql.Run(r.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	out := make([]Repository, 0)
	for _, n := range resp.ListNodes {
		r := &Repository{
			mgr:  r,
			name: n.Name,
			host: n.Host,
		}

		err = r.init()
		if err != nil {
			return nil, err
		}

		out = append(out, *r)
	}

	return out, nil
}

// Connect establishes a new repository connection.
func (r *RepositoriesManager) Connect(name, addr, key string, skipInsecureCheck bool) error {

	args := make([]string, 0)
	args = append(args, fmt.Sprintf(`name: "%s"`, name))
	args = append(args, fmt.Sprintf(`addr: "%s"`, addr))
	if key != "" {
		args = append(args, fmt.Sprintf(`credentials: "%s"`, key))
	}
	args = append(args, fmt.Sprintf(`insecureSkipVerify: %v`, skipInsecureCheck))
	argsStr := fmt.Sprintf("(%s)", strings.Join(args, ", "))

	req := graphql.NewRequest(fmt.Sprintf(`
		mutation {
			newNode%s{
				name
			}
		}
	`, argsStr))

	type responseContainer struct {
		NewNode objects.Node `json:"newNode"`
	}

	resp := new(responseContainer)
	err := r.c.graphql.Run(r.c.ctx, req, &resp)
	if err != nil {
		return err
	}

	return nil
}

// Get a specific repository.
func (r *RepositoriesManager) Get(name string) (*Repository, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
		query {
			listNodes {
				name
				host
			}
		}
	`))

	type responseContainer struct {
		ListNodes []objects.Node `json:"listNodes"`
	}
	resp := new(responseContainer)
	err := r.c.graphql.Run(r.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	for _, n := range resp.ListNodes {
		if n.Name == name {
			out := &Repository{
				mgr:  r,
				name: n.Name,
				host: n.Host,
			}
			err = out.init()
			if err != nil {
				return nil, err
			}

			return out, nil
		}
	}

	return nil, fmt.Errorf("could not repository '%s'", name)
}

// Disconnect destroys the Repository object and unregisters it from the current
// Vorteil environment.
func (r *RepositoriesManager) Disconnect(name string) error {

	req := graphql.NewRequest(fmt.Sprintf(`
		mutation {
			removeNode(name: "%s")
		}
	`, name))

	type responseContainer struct {
		RemoveNode bool `json:"removeNode"`
	}
	resp := new(responseContainer)
	err := r.c.graphql.Run(r.c.ctx, req, &resp)
	return err
}
