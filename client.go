package goapi

import (
	"context"
	"fmt"
	"strings"

	"code.vorteil.io/vorteil/apis/goapi/pkg/graphqlws"
	"github.com/machinebox/graphql"
)

// NewClient ..
func NewClient(cfg *ClientConfig) (*Client, error) {

	c := &Client{
		ctx:         context.Background(),
		cfg:         cfg,
		machinesMgr: &MachinesManager{},
		reposMgr: &RepositoriesManager{
			Local: &Repository{
				name: "local",
			},
		},
		buildMgr: &BuildManager{},
	}
	c.machinesMgr.c = c
	c.machinesMgr.environment = c.reposMgr.Local
	c.reposMgr.c = c
	c.buildMgr.c = c
	c.buildMgr.environment = c.reposMgr.Local

	var err error
	err = c.init()
	if err != nil {
		return nil, err
	}

	c.graphql = graphql.NewClient(fmt.Sprintf("%s%s/graphql", c.protocol, c.cfg.Address))
	c.subscriptions, err = graphqlws.NewClient(c.ctx, &graphqlws.ClientConfig{
		Address: c.cfg.Address,
		Path:    "subscriptions",
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Client provides access to the Vorteil API by establishing a connection to the
// specified Vorteil environment.
type Client struct {
	ctx           context.Context
	cfg           *ClientConfig
	protocol      string
	reposMgr      *RepositoriesManager
	machinesMgr   *MachinesManager
	buildMgr      *BuildManager
	subscriptions *graphqlws.Client
	graphql       *graphql.Client
}

// ClientConfig contains fields essential for the configuration of a new Client
type ClientConfig struct {
	Address           string
	AuthenticationKey string
}

func (c *Client) init() error {

	if c.cfg.Address == "" {
		return fmt.Errorf("address field may not be empty")
	}
	if strings.HasPrefix(c.cfg.Address, "https://") {
		c.protocol = "https://"
		c.cfg.Address = strings.TrimPrefix(c.cfg.Address, "https://")
	} else if strings.HasPrefix(c.cfg.Address, "http://") {
		c.cfg.Address = strings.TrimPrefix(c.cfg.Address, "http://")
	}
	if c.protocol == "" {
		c.protocol = "http://"
	}

	c.reposMgr.Local.mgr = c.reposMgr
	c.reposMgr.Local.hdr = make(map[string][]string)
	c.reposMgr.Local.host = fmt.Sprintf("%s%s", c.protocol, c.cfg.Address)

	return nil
}

// MachinesManager ..
func (c *Client) Machines() *MachinesManager {
	return c.machinesMgr
}

// BuildsManager ..
func (c *Client) Builds() *BuildManager {
	return c.buildMgr
}

// RepositoriesManager ..
func (c *Client) Repositories() *RepositoriesManager {
	return c.reposMgr
}
