package graphqlws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// TODO: keep-alive dependency
// TODO: GQL subscription timeout
// TODO: connect on-demand
// TODO: automatic reconnect
// TODO: accept-encoding gzip
// TODO: GQL_START extensions

// ClientConfig aggregates adjustable settings that can be used to modify the
// behaviour of the GraphQL web socket client.
type ClientConfig struct {

	// Logger can be used to customize the way the client logs information.
	// A nil value on this field will cause all logs to be discarded. Logs
	// are only meant to be used for troubleshooting, so this is usually the
	// desired behaviour.
	Logger Logger

	// Address defines the address component of the URL that will be used in
	// the request to connect to the server.
	Address string

	// Path defines the path component of the URL that will be used in the
	// request to connect to the server.
	Path string

	// Dialer can be used to provide the 'websocket.Dialer' that will be
	// used to connect to a server. This may be useful if the client needs
	// to be proxied or make use of a cookiejar. When left as nil, the
	// client will dial using the default dialer (websocket.DefaultDialer).
	Dialer *websocket.Dialer

	// Header can be used to provide a HTTP header that will be combined
	// with the header the client will use in its initial HTTP request.
	// This may be useful if the server handles authentication and
	// authorization via headers. For most cases it will be okay to leave
	// this as nil.
	Header http.Header

	// InitialPayload will be sent to the server during the first stages
	// initializing the connection. A GraphQL web socket server's
	// expectations regarding the initial payload can vary from
	// implementation to implementation, but it is a common requirement that
	// this payload contain information used to authenticate and authorize
	// the client.
	InitialPayload interface{}

	// ReadTimeout determines the maximum length of time a client will wait
	// in between each message received from the server. This setting only
	// comes into effect if the connected server is configured to
	// 'keep-alive' the connection. While in effect, if this interval is
	// ever exceeded the connection will be considered corrupt, causing the
	// client to close itself down. If left as zero, no read timeout will
	// be enforced.
	ReadTimeout time.Duration

	// WriteTimeout determines the maximum length of time a client will wait
	// for the server to receive each message. If this interval is ever
	// exceeded the connection will be considered corrupt, causing the
	// client to close itself down. If left as zero, no write timeout will
	// be enforced.
	WriteTimeout time.Duration
}

// Client manages a single connection to a GraphQL web socket server, which may
// host any number of GraphQL queries, mutations, and subscriptions
// concurrently. A Client object must be intialized using the NewClient
// function.
type Client struct {
	config          *ClientConfig
	log             logger
	url             url.URL
	conn            *websocket.Conn
	inbox           chan *Message
	readLoopClosed  chan bool
	expectKeepAlive bool

	// error propagation
	err         error
	errReported bool
	errLock     sync.RWMutex

	threads sync.WaitGroup
	cleanup []func()

	// outbox
	outbox       chan *Message
	outboxClosed bool
	outboxLock   sync.RWMutex

	// operations
	inShutdown bool
	operations map[string]*clientOperation
	opsLock    sync.RWMutex
}

// Subscription provides a handle for an active subscription with functions to
// terminate the subscription and to wait until the subscription has finished.
type Subscription struct {
	op *clientOperation
}

func (s *Subscription) Stop() {
	s.op.Stop()
}

func (s *Subscription) WaitUntilFinished(ctx context.Context) error {
	return s.op.WaitUntilFinished(ctx)
}

func (c *Client) dial(ctx context.Context) error {
	dialer := c.config.Dialer
	if dialer == nil {
		dialer = websocket.DefaultDialer
	}
	header := c.config.Header
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Sec-WebSocket-Protocol", "graphql-ws")
	var err error
	c.log.Info("Dialing server")
	c.conn, _, err = dialer.DialContext(ctx, c.url.String(), header)
	if err != nil {
		return err
	}
	c.log.Info("Connected to server")
	c.cleanup = append([]func(){func() {
		c.log.Info("Closing the connection to the server")
		c.reportError(c.conn.Close())
	}}, c.cleanup...)

	protocol := c.conn.Subprotocol()
	if protocol != "graphql-ws" {
		return errors.New("failed to negotiate 'graphql-ws' subprotocol")
	}
	return nil
}

func (c *Client) reportError(err error) {
	c.errLock.Lock()
	if !c.errReported {
		c.errReported = true
		c.err = err
	}
	c.errLock.Unlock()
}

func (c *Client) newWriteDeadline() time.Time {
	var t time.Time
	if c.config.WriteTimeout != 0 {
		t = time.Now().Add(c.config.WriteTimeout)
	}
	return t
}

func (c *Client) writeLoop() {
	c.log.Info("Starting write loop")
	for {
		msg, more := <-c.outbox
		if !more {
			c.log.Info("Write loop terminating because the outbox has been closed")
			break
		}

		err := c.conn.SetWriteDeadline(c.newWriteDeadline())
		if err != nil {
			err = fmt.Errorf("failed to set a write deadline: %v", err)
			c.reportError(err)
			break
		}

		err = c.conn.WriteJSON(msg)
		if err != nil {
			err = fmt.Errorf("failed to write to the connection: %v", err)
			c.reportError(err)
			break
		}
	}
	c.log.Info("Write loop finished")
	c.threads.Done()
}

func (c *Client) startWriteLoop() {
	c.threads.Add(1)
	c.cleanup = append([]func(){func() {
		c.outboxLock.Lock()
		if !c.outboxClosed {
			close(c.outbox)
			c.outboxClosed = true
		}
		c.outboxLock.Unlock()
	}}, c.cleanup...)
	go c.writeLoop()
}

func (c *Client) newReadDeadline() time.Time {
	var t time.Time
	if c.expectKeepAlive {
		if c.config.ReadTimeout != 0 {
			t = time.Now().Add(c.config.ReadTimeout)
		}
	}
	return t
}

func (c *Client) readLoop() {
	c.log.Info("Starting read loop")
	for {
		err := c.conn.SetReadDeadline(c.newReadDeadline())
		if err != nil {
			err = fmt.Errorf("failed to set a read deadline: %v", err)
			c.reportError(err)
			break
		}

		msg := new(Message)
		err = c.conn.ReadJSON(&msg)
		if err != nil {
			err = fmt.Errorf("failed to read from the connection: %v", err)
			c.reportError(err)
			break
		}

		c.log.Info("Read loop queuing a new message for the dispatcher")
		c.inbox <- msg
	}

	close(c.inbox)
	c.log.Info("Read loop finished")
	c.threads.Done()

	// force outbox to close in case the closed connection was not intentional
	close(c.readLoopClosed)
	c.outboxLock.Lock()
	if !c.outboxClosed {
		close(c.outbox)
		c.outboxClosed = true
	}
	c.outboxLock.Unlock()
}

func (c *Client) startReadLoop() {
	c.threads.Add(1)
	c.cleanup = append([]func(){func() {
		// reporting a nil error causes all subsequent errors to be disregarded
		c.reportError(nil)
		err := c.conn.Close()
		if err != nil {
			c.reportError(err)
		}
	}}, c.cleanup...)
	go c.readLoop()
}

func (c *Client) dispatcher() {
	c.log.Info("Starting dispatcher loop")
	for {
		msg, more := <-c.inbox
		if !more {
			break
		}

		c.log.Info("Dispatcher received a new message")

		switch msg.Type {
		case MessageTypeGQLConnectionKeepAlive:
			c.log.Info("Received keep-alive message from server")
			c.expectKeepAlive = true
		case MessageTypeGQLData:
			c.opsLock.RLock()
			o, ok := c.operations[msg.ID]
			if ok {
				data, err := msg.GQLDataPayload()
				if err != nil {
					o.log.Error("Discarding corrupt data payload: %v", err)
				} else {
					o.data(data)
				}
			} else {
				c.log.Error("Discarding data payload for unknown operation: %s", msg.ID)
			}
			c.opsLock.RUnlock()
		case MessageTypeGQLError:
			c.opsLock.RLock()
			o, ok := c.operations[msg.ID]
			if ok {
				o.error(msg.Payload)
			} else {
				c.log.Error("Discarding error for an unknown operation: %s: %v", msg.ID, msg.Payload)
			}
			c.opsLock.RUnlock()
		case MessageTypeGQLComplete:
			c.opsLock.Lock()
			o, ok := c.operations[msg.ID]
			if ok {
				delete(c.operations, msg.ID)
				o.complete()
				o.log.Info("Cleaned up completed operation")
			} else {
				c.log.Error("Server indicated an unknown operation was completed: %s", msg.ID)
			}
			c.opsLock.Unlock()
		case MessageTypeGQLConnectionError:
			c.log.Error("Server ignored a message due to parsing errors: %v", msg.Payload)
		default:
			c.log.Error("Server sent unexpected message type: %s", msg.Type)
		}
	}

	c.log.Info("Dispatcher loop finished")
	c.threads.Done()
}

func (c *Client) startDispatcher() {
	c.threads.Add(1)
	c.cleanup = append([]func(){func() {}}, c.cleanup...)
	go c.dispatcher()
}

func (c *Client) send(message *Message) {
	c.outboxLock.RLock()
	if !c.outboxClosed {
		c.outbox <- message
	}
	c.outboxLock.RUnlock()
}

func (c *Client) initConnection(ctx context.Context) error {
	c.startWriteLoop()
	c.startReadLoop()
	c.log.Info("Initializing websocket with graphql-ws subprotocol")
	// send GraphQL connection init message
	if c.config.InitialPayload == nil {
		c.config.InitialPayload = make(map[string]interface{})
	}
	c.log.Info("Queuing GraphQL connection init message")
	c.send(&Message{
		Type:    MessageTypeGQLConnectionInit,
		Payload: c.config.InitialPayload,
	})
	c.cleanup = append([]func(){
		func() {
			c.log.Info("Queuing GraphQL connection terminate message")
			c.send(&Message{
				Type:    MessageTypeGQLConnectionTerminate,
				Payload: map[string]interface{}{},
			})
		},
	}, c.cleanup...)

	// wait for server's acknowledgement of the connection
	var m *Message
	select {
	case msg, more := <-c.inbox:
		if !more {
			return fmt.Errorf("connection closed: %v", c.err)
		}
		m = msg
	case <-ctx.Done():
		return fmt.Errorf("context cancelled or timed out: %v", ctx.Err())
	}
	// check that the response matches our expectations
	switch m.Type {
	case MessageTypeGQLConnectionAck:
	case MessageTypeGQLConnectionError:
		return fmt.Errorf("the server rejected the connection: %v", m.Payload)
	default:
		return fmt.Errorf("server responded to GraphQL connection init message with unexpected message type: %s", m.Type)
	}

	c.startDispatcher()

	// defer a cleanup of each clientOperation
	c.cleanup = append([]func(){
		func() {
			var liveOperations []*clientOperation
			c.opsLock.RLock()
			for _, v := range c.operations {
				liveOperations = append(liveOperations, v)
			}
			c.opsLock.RUnlock()

			var wg sync.WaitGroup
			wg.Add(len(liveOperations))
			for _, o := range liveOperations {
				go func(op *clientOperation) {
					op.Stop()
					_ = op.WaitUntilFinished(context.Background())
					wg.Done()
				}(o)
			}
			wg.Wait()
			c.log.Info("Stopped all live operations")
		},
	}, c.cleanup...)
	return nil
}

// NewClient creates and connects a new GraphQL web socket client to a GraphQL
// web socket server.
func NewClient(ctx context.Context, config *ClientConfig) (*Client, error) {
	c := new(Client)
	c.config = config
	c.url = url.URL{Scheme: "ws", Host: config.Address, Path: config.Path}
	// initialize logger
	c.log.logger = c.config.Logger

	c.readLoopClosed = make(chan bool)
	c.inbox = make(chan *Message)
	c.outbox = make(chan *Message)
	c.operations = make(map[string]*clientOperation)

	// dial
	err := c.dial(ctx)
	if err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("failed to dial server: %v", err)
	}

	// initialize connection
	err = c.initConnection(ctx)
	if err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("failed to initialize connection: %v", err)
	}
	return c, nil
}

// Shutdown attempts to close the client gracefully, terminating all ongoing
// operations correctly before closing the websocket and terminating the
// connection. The provided context can be used to cancel the shutdown
// prematurely in the event that it stalls or takes longer than is acceptable.
func (c *Client) Shutdown(ctx context.Context) error {
	ch := make(chan bool)
	var cancelled bool
	c.opsLock.Lock()
	if c.inShutdown {
		c.opsLock.Unlock()
		return errors.New("client is already in shutdown")
	}
	c.inShutdown = true
	c.opsLock.Unlock()

	c.log.Info("Shutting down client")

	go func() {
		for _, fn := range c.cleanup {
			if cancelled {
				break
			}
			fn()
		}
		c.threads.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		c.log.Info("Client has shut down")
		return nil
	case <-ctx.Done():
		cancelled = true
		return fmt.Errorf("failed to shutdown the client: %v", ctx.Err())
	}
}

// Close is a non-blocking function that immediately terminates the connection
// and cleans up resources in use by the Client.
func (c *Client) Close() error {
	c.log.Info("Closing client")
	c.inShutdown = true
	var wg sync.WaitGroup
	wg.Add(len(c.cleanup))
	for _, fn := range c.cleanup {
		go func(fn func()) {
			fn()
			wg.Done()
		}(fn)
	}
	go func() {
		wg.Wait()
		c.log.Info("Client closed")
	}()
	return nil
}

type queryConfig struct {
	Query         string
	Variables     map[string]interface{}
	OperationName string
}

func (c *queryConfig) validateIsType(t string) error {
	document, err := parser.Parse(parser.ParseParams{
		Source: c.Query,
	})
	if err != nil {
		return fmt.Errorf("failed to parse query as valid GraphQL: %v", err)
	}

	if l := len(document.Definitions); l == 0 {
		return errors.New("query string must define a valid GraphQL operation")
	}
	var index = -1
	for i, definition := range document.Definitions {
		if definition.GetKind() != "OperationDefinition" {
			return errors.New("non OperationDefinition in document")
		}
		if index >= 0 {
			return errors.New("query string cannot define multiple operations")
		}
		index = i
	}
	if index < 0 {
		return errors.New("query string must define a valid GraphQL operation")
	}
	definition := document.Definitions[index].(*ast.OperationDefinition)
	if definition.Operation != t {
		return fmt.Errorf("query string must define a %s, instead found '%s'", t, definition.Operation)
	}

	return nil
}

type operationConfig struct {
	Query         string
	Variables     map[string]interface{}
	OperationName string
	DataCallback  func(payload *GQLDataPayload)
	ErrorCallback func(err error)
}

type clientOperation struct {
	*Client
	log      operationLogger
	id       string
	cfg      operationConfig
	finished chan bool
}

func (o *clientOperation) complete() {
	delete(o.operations, o.id)
	close(o.finished)
}

func (o *clientOperation) error(x interface{}) {
	o.log.Info("Delivering GraphQL error message payload via ErrorCallback")
	s := fmt.Sprintf("%v", x)
	err := errors.New(s)
	o.cfg.ErrorCallback(err)
}

func (o *clientOperation) data(payload *GQLDataPayload) {
	o.log.Info("Delivering GraphQL data message payload via DataCallback")
	o.cfg.DataCallback(payload)
}

func (o *clientOperation) Stop() {
	select {
	case <-o.finished:
	default:
		o.log.Info("Queuing GraphQL stop message")
		o.send(&Message{
			ID:   o.id,
			Type: MessageTypeGQLStop,
		})
	}
}

func (o *clientOperation) WaitUntilFinished(ctx context.Context) error {
	select {
	case <-o.finished:
	case <-o.readLoopClosed:
		o.opsLock.Lock()
		o.complete()
		o.opsLock.Unlock()
		o.log.Info("Cleaned up operation because of a closed connection")
	case <-ctx.Done():
		return fmt.Errorf("context cancelled or timed out before client finished: %v", ctx.Err())
	}
	return nil
}

func (c *Client) beginOperation(config operationConfig) (*clientOperation, error) {
	o := new(clientOperation)
	o.Client = c
	o.cfg = config
	o.finished = make(chan bool)

	c.opsLock.Lock()
	if c.inShutdown {
		c.opsLock.Unlock()
		return nil, errors.New("cannot begin a new operation while client is in shutdown")
	}

	for {
		uid, err := uuid.NewV4()
		if err != nil {
			c.opsLock.Unlock()
			return nil, fmt.Errorf("failed to generate uuid: %v", err)
		}
		o.id = uid.String()
		if _, exists := c.operations[o.id]; !exists {
			break
		}
	}

	o.log.logger = o.Client.log
	o.log.suffix = fmt.Sprintf(" (operation: %s)", o.id)
	c.operations[o.id] = o
	c.opsLock.Unlock()

	o.log.Info("Added a new operation to the client")

	go o.initialize()

	return o, nil
}

func (o *clientOperation) initialize() {
	o.log.Info("Queuing GraphQL start message")
	o.send(&Message{
		ID:   o.id,
		Type: MessageTypeGQLStart,
		Payload: GQLStartPayload{
			Query:         o.cfg.Query,
			Variables:     o.cfg.Variables,
			OperationName: o.cfg.OperationName,
		},
	})
}

type onceConfig interface {
	validate() error
	queryConfig() queryConfig
}

func (c *Client) once(ctx context.Context, config onceConfig) (*GQLDataPayload, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}

	var result *GQLDataPayload
	var gqlError error
	cfg := config.queryConfig()

	o, err := c.beginOperation(operationConfig{
		Query:         cfg.Query,
		Variables:     cfg.Variables,
		OperationName: cfg.OperationName,
		DataCallback: func(pl *GQLDataPayload) {
			result = pl
		},
		ErrorCallback: func(e error) {
			gqlError = e
		},
	})
	if err != nil {
		return nil, err
	}

	err = o.WaitUntilFinished(ctx)
	if err != nil {
		return nil, err
	}
	if gqlError != nil {
		return nil, gqlError
	}

	return result, nil
}

// QueryConfig contains fields used by the client to perform a GraphQL query
// clientOperation.
type QueryConfig struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

func (c *QueryConfig) validate() error {
	return (*queryConfig)(c).validateIsType("query")
}

func (c QueryConfig) queryConfig() queryConfig {
	return queryConfig(c)
}

// Query performs a GraphQL query over the client's websocket.
func (c *Client) Query(ctx context.Context, config *QueryConfig) (*GQLDataPayload, error) {
	return c.once(ctx, config)
}

// MutationConfig contains fields used by the client to perform a GraphQL
// mutation clientOperation.
type MutationConfig struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

func (c *MutationConfig) validate() error {
	return (*queryConfig)(c).validateIsType("mutation")
}

func (c *MutationConfig) queryConfig() queryConfig {
	return queryConfig(*c)
}

// Mutation performs a GraphQL mutation over the client's websocket.
func (c *Client) Mutation(ctx context.Context, config *MutationConfig) (*GQLDataPayload, error) {
	return c.once(ctx, config)
}

// SubscriptionConfig contains fields used by the client to perform a GraphQL
// subscription clientOperation.
type SubscriptionConfig struct {
	Query         string
	Variables     map[string]interface{}
	OperationName string
	DataCallback  func(payload *GQLDataPayload)
	ErrorCallback func(err error)
}

func (c *SubscriptionConfig) validate() error {
	return (&queryConfig{
		Query:         c.Query,
		Variables:     c.Variables,
		OperationName: c.OperationName,
	}).validateIsType("subscription")
}

// Subscription registers a GraphQL subscription for repeated updates over the
// client's websocket.
func (c *Client) Subscription(config *SubscriptionConfig) (*Subscription, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}

	o, err := c.beginOperation(operationConfig(*config))
	if err != nil {
		return nil, err
	}

	s := new(Subscription)
	s.op = o

	return s, nil
}
