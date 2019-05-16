package goapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/machinebox/graphql"
)

type Scheme string

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client ..
type Client struct {
	http   *http.Client
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
		cfg:  cfg,
		ctx:  context.Background(),
		http: http.DefaultClient,
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

// CursorArgs ..
type CursorArgs struct {
	After  string `json:"after"`
	Before string `json:"before"`
	First  int    `json:"first"`
	Last   int    `json:"last"`
}

// ParseCursor ..
func parseCursor(curs *CursorArgs) (string, string) {
	var variableDeclarations string
	var variables string
	var cursorPresent bool

	if curs != nil {
		if curs.After != "" {
			cursorPresent = true
			if len(variableDeclarations) != 0 {
				variableDeclarations = fmt.Sprintf("%s,", variableDeclarations)
			}
			if len(variables) != 0 {
				variables = fmt.Sprintf("%s,", variables)
			}
			variables = fmt.Sprintf("%safter:$after", variables)
			variableDeclarations = fmt.Sprintf("%s $after: String", variableDeclarations)
		}
		if curs.Before != "" {
			cursorPresent = true
			if len(variableDeclarations) != 0 {
				variableDeclarations = fmt.Sprintf("%s,", variableDeclarations)
			}
			if len(variables) != 0 {
				variables = fmt.Sprintf("%s,", variables)
			}
			variables = fmt.Sprintf("%sbefore:$before", variables)
			variableDeclarations = fmt.Sprintf("%s $before: String", variableDeclarations)
		}
		if curs.First != 0 {
			cursorPresent = true
			if len(variableDeclarations) != 0 {
				variableDeclarations = fmt.Sprintf("%s,", variableDeclarations)
			}
			if len(variables) != 0 {
				variables = fmt.Sprintf("%s,", variables)
			}
			variables = fmt.Sprintf("%sfirst:$first", variables)
			variableDeclarations = fmt.Sprintf("%s$first: Int", variableDeclarations)
		}
		if curs.Last != 0 {
			cursorPresent = true
			if len(variableDeclarations) != 0 {
				variableDeclarations = fmt.Sprintf("%s,", variableDeclarations)
			}
			if len(variables) != 0 {
				variables = fmt.Sprintf("%s,", variables)
			}
			variables = fmt.Sprintf("%slast:$last", variables)
			variableDeclarations = fmt.Sprintf("%s $last: Int", variableDeclarations)
		}

		if cursorPresent {
			variables = fmt.Sprintf("(%s)", variables)
			variableDeclarations = fmt.Sprintf("(%s)", variableDeclarations)
		}
	}

	return variableDeclarations, variables
}

func addCursorToRequest(req *graphql.Request, curs *CursorArgs) {
	if curs.After != "" {
		req.Var("after", curs.After)
	}
	if curs.Before != "" {
		req.Var("before", curs.Before)
	}
	if curs.First != 0 {
		req.Var("first", curs.First)
	}
	if curs.Last != 0 {
		req.Var("last", curs.Last)
	}
}

// Do ..
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	// r.Header.Set(VersionHeaderKey, Version.Release)

	// if c.Forward != "" {
	// 	c.Logger.Debug(fmt.Sprintf("Forwarding request: %s", c.Forward))
	// 	r.Header.Set("Vorteil", c.Forward)
	// }
	resp, err := c.http.Do(r)
	if err != nil {
		return nil, err
	}

	// v := resp.Header.Get(VersionHeaderKey)
	// x := strings.SplitN(v, ".", 3)

	// var maj, min uint
	// maj = uint(ContractMajor)
	// min = uint(ContractMinor)

	// cont := contracts.Contract{
	// 	Major: &maj,
	// 	Minor: &min,
	// }

	return resp, nil

}
