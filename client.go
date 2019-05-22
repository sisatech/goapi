package goapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/inconshreveable/log15"
	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/graphqlws"
)

// Scheme ...
type Scheme string

// Doer ...
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client ..
type Client struct {
	http          *http.Client
	cfg           *ClientConfig
	cookie        *http.Cookie
	client        *graphql.Client
	ctx           context.Context
	subscriptions *graphqlws.Client
	basicAuth     *basicAuth
	logger        *logger
}

type basicAuth struct {
	headerKey string
	headerVal string
}

// ClientConfig - used to create a new Client
type ClientConfig struct {
	Address string
	Path    string
	WSPath  string
	Key     string
	Logger  log15.Logger
}

type logger struct {
	log log15.Logger
}

func NewLogger(log log15.Logger) *logger {
	if log == nil {
		log = log15.Root()
	}

	return &logger{log: log}
}

func (l *logger) Info(message string) {
	l.log.Info(message)
}

func (l *logger) Error(message string) {
	l.log.Error(message)
}

// NewClient returns a Client according to the provided *ClientArgs
func NewClient(ctx context.Context, cfg *ClientConfig) (*Client, error) {

	var err error

	l := NewLogger(cfg.Logger)
	c := &Client{
		cfg:    cfg,
		ctx:    ctx,
		http:   http.DefaultClient,
		logger: l,
	}

	if cfg.Key != "" {
		err = c.KeyAuthentication(cfg.Key)
		if err != nil {
			return nil, err
		}
	}

	clientURL := fmt.Sprintf("http://%s", c.graphqlURL())
	c.client = graphql.NewClient(clientURL)

	c.subscriptions, err = graphqlws.NewClient(c.ctx, &graphqlws.ClientConfig{
		Address: c.cfg.Address,
		Path:    c.cfg.WSPath,
		Logger:  c.logger,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

// KeyAuthentication ..
func (c *Client) KeyAuthentication(key string) error {

	resp, err := c.http.Post(fmt.Sprintf("http://%s/api/login", c.cfg.Address),
		"application/json", bytes.NewReader([]byte(fmt.Sprintf(`{"key": "%s"}`, c.cfg.Key))))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response had a non-200 status code: %v", resp.StatusCode)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "vauth" {
			key = c.Value
		}
	}

	c.basicAuth = &basicAuth{
		headerKey: "Cookie",
		headerVal: "vauth=" + key,
	}

	header := make(http.Header)
	header.Set(c.basicAuth.headerKey, c.basicAuth.headerVal)

	c.subscriptions, err = graphqlws.NewClient(c.ctx, &graphqlws.ClientConfig{
		Address: c.cfg.Address,
		Path:    c.cfg.WSPath,
		Header:  header,
	})
	if err != nil {
		return err
	}

	return nil
}

// BasicAuthentication ..
func (c *Client) BasicAuthentication(user, pw string) error {

	creds := []byte(fmt.Sprintf("%s:%s", user, pw))
	credStr := "Basic " + base64.StdEncoding.EncodeToString(creds)

	header := make(http.Header)
	header.Set("Authorization", credStr)

	c.basicAuth = &basicAuth{
		headerKey: "Authorization",
		headerVal: credStr,
	}

	var err error
	c.subscriptions, err = graphqlws.NewClient(c.ctx, &graphqlws.ClientConfig{
		Address: c.cfg.Address,
		Path:    c.cfg.WSPath,
		Logger:  c.logger,
		Header:  header,
	})
	if err != nil {
		return err
	}

	return nil
}

// NewRequest ...
func (c *Client) NewRequest(str string) *graphql.Request {
	req := graphql.NewRequest(str)
	if c.basicAuth != nil {
		req.Header.Set(c.basicAuth.headerKey, c.basicAuth.headerVal)
	}

	return req
}

// SetContext ...
func (c *Client) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// Context ...
func (c *Client) Context() context.Context {
	return c.ctx
}

// Cookie ...
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

	base := c.cfg.Address
	if !strings.HasPrefix(base, "https://") && !strings.HasPrefix(base, "http://") {
		base = fmt.Sprintf("http://%s", base)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", base, url), body)
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

	base := c.cfg.Address
	if !strings.HasPrefix(base, "https://") && !strings.HasPrefix(base, "http://") {
		base = fmt.Sprintf("http://%s", base)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", base, url), nil)
	if err != nil {
		return nil, err
	}
	if c.Cookie() != nil {
		req.AddCookie(c.cookie)
	}
	return http.DefaultClient.Do(req)
}

// func (c *Client) NewRequest(str string) *graphql.Request {
// 	req := graphql.NewRequest(str)
// 	req.Header.Set("")
// }

// Do ..
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	// r.Header.Set(VersionHeaderKey, Version.Release)

	// if c.Forward != "" {
	// 	c.Logger.Debug(fmt.Sprintf("Forwarding request: %s", c.Forward))
	// 	r.Header.Set("Vorteil", c.Forward)
	// }
	// resp, err := c.http.Do(r)
	// if err != nil {
	// 	return nil, err
	// }

	// v := resp.Header.Get(VersionHeaderKey)
	// x := strings.SplitN(v, ".", 3)

	// var maj, min uint
	// maj = uint(ContractMajor)
	// min = uint(ContractMinor)

	// cont := contracts.Contract{
	// 	Major: &maj,
	// 	Minor: &min,
	// }

	return c.http.Do(r)

}
