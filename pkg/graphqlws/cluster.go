package graphqlws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Peer represents another node within a cluster. It is used by implementation
// of Cluster to notify other nodes in the cluster that an Update has been
// published.
type Peer interface {
	// Notify sends an Update to the peer.
	Notify(update Update) error
}

// Peers represents the collection of all peers. It exists as an interface in
// case the implementation wants to handle some management of the underlying
// peers.
type Peers interface {
	// Peers returns a complete list of Peer objects in the cluster.
	Peers() []Peer
}

// SimplePeer is the basic implementation of the Peer interface.
type SimplePeer struct {
	Client *http.Client
	Header http.Header
	URL    *url.URL
}

// Notify implements the Peer interface by sending a POST request to the
// SimplePeer's URL. The POST request conforms to the expectations of the
// Cluster's ServeHTTP handler so that the peer will trigger an update of its
// GraphQL Web Socket server.
func (p *SimplePeer) Notify(update Update) error {
	data, err := json.MarshalIndent(update, "", "  ")
	if err != nil {
		return err
	}
	body := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, p.URL.String(), body)
	if err != nil {
		return err
	}
	if p.Header != nil {
		for key, v := range p.Header {
			for i, val := range v {
				if i == 0 {
					req.Header.Set(key, val)
				} else {
					req.Header.Add(key, val)
				}
			}
		}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("peer responded to notification with non-200 status code: %v", resp.StatusCode)
	}
	return nil
}

// StaticPeers makes a type out of a slice of Peer object so methods can be
// attached in order to satisfy the Peers interface.
type StaticPeers []Peer

// Peers satisfies the Peers interface.
func (ps StaticPeers) Peers() []Peer {
	return ps
}

// Clusterable is used as a part of the ClusterConfig in order to create an
// implementation of the Cluster interface based on any underlying Publisher.
type Clusterable interface {

	// Something Clusterable must be capable of publishing updates.
	Publish(update Update)

	// Schema returns an implementation of Schema, which is used by the
	// implementation of Cluster to reconstruct Update objects that have
	// been received from peers.
	Schema() Schema
}

// TODO: peer auth config

// ClusterConfig contains fields used to customize the behaviour of a basic
// implementation of the Cluster interface.
type ClusterConfig struct {

	// Clusterable is the Clusterable object to be managed by the Cluster.
	Clusterable Clusterable

	// Logger can be used to configure the log output of the Cluster.
	Logger Logger

	// Is the Peers object that is used to communicate with the other nodes
	// of the cluster.
	Peers Peers
}

// Cluster represents an entire cluster of Publishers working as a HA
// (highly-available) unit. By clustering Publishers load can be distributed
// across multiple machines, and implementations of the Cluster interface
// ensure that updates published to one node also reach subscribers to the other
// nodes.
type Cluster struct {
	core             Clusterable
	log              logger
	maxPayloadLength int
	peers            Peers
}

// NewCluster returns a basic implementation of the Cluster interface assembled
// from implementations of the various critical elements of the provided
// ClusterConfig.
//
// Although it notifies peers via HTTP, it does not perform any complicated
// logic to guarantee that all peers receive their notification. If an error
// occurs in an attempt to notify a peer, the error is logged and ignored.
func NewCluster(config ClusterConfig) (*Cluster, error) {
	c := new(Cluster)
	c.core = config.Clusterable
	c.peers = config.Peers
	c.log.logger = config.Logger

	// marshal entire scheme as false to estimate an upper bound to update payload size
	u := c.core.Schema().NewUpdate()
	pl, err := u.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to estimate maximum update payload size: %v", err)
	}
	c.maxPayloadLength = 1024 * len(pl)

	return c, nil
}

// Publish pushed out an update to all nodes in the cluster. It is an
// implementation of Publisher, making it easy to swap in and out with
// non-clustered Publisher implementations.
func (c *Cluster) Publish(update Update) {

	c.log.Info("Publishing update to cluster")

	go c.core.Publish(update)

	peers := c.peers.Peers()
	for _, peer := range peers {
		go func(p Peer) {
			c.log.Info("Pushing update to a peer")
			err := p.Notify(update)
			if err != nil {
				c.log.Error(err.Error())
			}
			c.log.Info("Pushed an update to a peer")
		}(peer)
	}

}

// ServeHTTP implements 'http.Handler' and should be served on a 'http.Server'
// somewhere. This handler is the input exposed by the Cluster that should be
// contacted by other nodes in the cluster in order to trigger an update.
//
// The handler expects a POST request with a payload containing JSON that will
// unmarshal into a valid Update description.
func (c *Cluster) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		msg := http.StatusText(code)
		http.Error(w, msg, code)
		return
	}

	if vals := r.URL.Query(); len(vals) > 0 {
		code := http.StatusBadRequest
		msg := "unacceptable url query parameters"
		http.Error(w, msg, code)
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		code := http.StatusBadRequest
		msg := "'Content-Type' header must be set to 'application/json'"
		http.Error(w, msg, code)
		return
	}

	if r.ContentLength == 0 {
		code := http.StatusBadRequest
		msg := "'Content-Length' header must be set and cannot be zero"
		http.Error(w, msg, code)
		return
	}

	if r.ContentLength > int64(c.maxPayloadLength) {
		code := http.StatusBadRequest
		msg := "request payload rejected for being oversize"
		http.Error(w, msg, code)
		return
	}

	buf := new(bytes.Buffer)
	k, err := io.CopyN(buf, r.Body, r.ContentLength)
	if err != nil && err != io.EOF {
		code := http.StatusBadRequest
		msg := "request payload rejected for being oversize"
		http.Error(w, msg, code)
		return
	}
	if err == io.EOF || k != r.ContentLength {
		code := http.StatusBadRequest
		msg := "'Content-Length' header doesn't match the length of the request body"
		http.Error(w, msg, code)
		return
	}

	data := buf.Bytes()
	update, err := c.core.Schema().UpdateFromJSON(data)
	if err != nil {
		code := http.StatusBadRequest
		msg := err.Error()
		http.Error(w, msg, code)
		return
	}

	c.log.Info("Received update from a peer")
	go c.core.Publish(update)

}
