package graphqlws

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
)

// TODO: PreventRepeats bool
// TODO: ContextInjector

type ServerConfig struct {
	Logger          Logger
	Schema          graphql.Schema
	EnableKeepAlive bool
	PollingInterval time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

type connectionManager struct {
	*Server
	lock   sync.RWMutex
	closed bool
	conns  map[string]*connection
}

func (s *Server) newConnectionManager() *connectionManager {
	return &connectionManager{
		Server: s,
		conns:  make(map[string]*connection),
	}
}

func (m *connectionManager) add(c *connection) {
	for {
		uid, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		c.id = uid.String()
		m.lock.Lock()
		if m.closed {
			return // TODO: return error
		}
		_, exists := m.conns[c.id]
		if exists {
			m.lock.Unlock()
			continue
		}
		m.conns[c.id] = c
		m.lock.Unlock()
		break
	}
}

func (m *connectionManager) remove(c *connection) {
	m.lock.Lock()
	delete(m.conns, c.id)
	m.lock.Unlock()
}

func (m *connectionManager) Range(fn func(c *connection)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, v := range m.conns {
		fn(v)
	}
}

type Server struct {
	cfg         ServerConfig
	log         logger
	schema      Schema
	upgrader    websocket.Upgrader
	connections *connectionManager
	updateLock  sync.Mutex
}

func NewServer(config *ServerConfig) (*Server, error) {
	s := new(Server)
	s.cfg = *config
	s.schema = NewSchema(s.cfg.Schema)
	s.log.logger = s.cfg.Logger
	s.upgrader = websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: []string{"graphql-ws"},
	}
	s.connections = s.newConnectionManager()
	return s, nil
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

func (s *Server) Publish(update Update) {
	s.updateLock.Lock()
	defer s.updateLock.Unlock()
	s.connections.Range(func(c *connection) {
		c.operations.Range(func(o *operation) {
			if !o.update.Overlaps(update) {
				return
			}
			o.execute()
		})
	})
}

func (s *Server) Schema() Schema {
	return s.schema
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Serving request")
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Info(fmt.Sprintf("Failed to establish websocket connection: %v", err))
		return
	}
	defer ws.Close()
	defer func() {
		_ = ws.WriteMessage(websocket.CloseMessage, []byte{})
	}()
	if ws.Subprotocol() != "graphql-ws" {
		s.log.Info("Connection does not implement the GraphQL WS protocol")
		return
	}
	connection := s.newConnection(ws)
	connection.run()
}
