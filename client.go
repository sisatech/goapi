package goapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/machinebox/graphql"
)

type Scheme string

// Client ..
type Client struct {
	cfg    *ClientConfig
	cookie *http.Cookie
	client *graphql.Client
	ctx    context.Context
}

// ClientArgs - used to create a new Client
type ClientConfig struct {
	Address string
	Path    string
}

// NewClient returns a Client according to the provided *ClientArgs
func NewClient(ctx context.Context, cfg *ClientConfig) (*Client, error) {

	c := &Client{
		cfg: cfg,
		ctx: context.Background(),
	}

	c.client = graphql.NewClient(c.graphqlURL())
	return c, nil
}

func (c *Client) NewRequest(str string) *graphql.Request {
	req := graphql.NewRequest(str)
	return req
}

func (c *Client) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *Client) Context() context.Context {
	return c.ctx
}

func (c *Client) Cookie() *http.Cookie {
	if c.cookie == nil {
		return nil
	}
	return c.cookie
}

func (c *Client) graphqlURL() string {
	return fmt.Sprintf("%s/graphql", c.cfg.Address)
}

func (c *Client) loginURL() string {
	return fmt.Sprintf("%s/api/login", c.cfg.Address)
}

// Post wraps http.Post()
func (c *Client) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	if c.Cookie() != nil {
		req.AddCookie(c.cookie)
	}
	return http.DefaultClient.Do(req)
}

// Get wraps http.Get()
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.Cookie() != nil {
		req.AddCookie(c.cookie)
	}
	return http.DefaultClient.Do(req)
}
