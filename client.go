package goapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/machinebox/graphql"
)

type Scheme string

const (
	HTTP  = Scheme("http")
	HTTPS = Scheme("https")
)

// Client ..
type Client struct {
	cfg    *ClientArgs
	cookie *http.Cookie
	client *graphql.Client
	ctx    context.Context
}

// ClientArgs - used to create a new Client
type ClientArgs struct {
	Host    string
	Port    int
	AuthKey string
	Scheme  Scheme
}

// NewTestClient returns a Client according to the provided *ClientArgs
func NewTestClient(args *ClientArgs) (*Client, error) {

	c := &Client{
		cfg: args,
		ctx: context.Background(),
	}

	c.client = graphql.NewClient(c.graphqlURL())
	if c.cfg.AuthKey != "" {
		// authenticate (get cookie)
		err := c.Login()
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) NewRequest(str string) *graphql.Request {
	req := graphql.NewRequest(str)
	if c.cfg.AuthKey != "" {
		req.Header.Set("Cookie", fmt.Sprintf("vauth=%s", c.cfg.AuthKey))
	}
	return req
}

func (c *Client) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *Client) Context() context.Context {
	return c.ctx
}

// Login with the AuthKey field from Client.cfg
func (c *Client) Login() error {

	resp, err := c.Post(c.loginURL(), mime.TypeByExtension("json"),
		bytes.NewReader([]byte(fmt.Sprintf("{\"key\": \"%s\"}",
			c.cfg.AuthKey))))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %v, expected %v",
			resp.StatusCode, http.StatusOK)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "vauth" {
			c.cookie = cookie
			return nil
		}
	}

	return fmt.Errorf("response did not contain 'vauth' cookie")
}

func (c *Client) Cookie() *http.Cookie {
	if c.cookie == nil {
		return nil
	}
	return c.cookie
}

func (c *Client) baseURL() string {
	return fmt.Sprintf("%s://%s:%v", c.cfg.Scheme, c.cfg.Host, c.cfg.Port)
}

func (c *Client) graphqlURL() string {
	return fmt.Sprintf("%s/graphql", c.baseURL())
}

func (c *Client) loginURL() string {
	return fmt.Sprintf("%s/api/login", c.baseURL())
}

// GraphQL ..
func (c *Client) GraphQL() *graphql.Client {
	return c.client
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
