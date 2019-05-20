package graphqlws

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type operationManager struct {
	*connection
	lock   sync.RWMutex
	closed bool
	ops    map[string]*operation
}

func (c *connection) newOperationManager() *operationManager {
	return &operationManager{
		connection: c,
		ops:        make(map[string]*operation),
	}
}

func (m *operationManager) get(uid string) (*operation, bool) {
	m.lock.RLock()
	if m.closed {
		m.lock.RUnlock()
		return nil, false
	}
	v, ok := m.ops[uid]
	m.lock.RUnlock()
	return v, ok
}

func (m *operationManager) add(o *operation) {
	for {
		uid, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		o.uid = uid.String()
		m.lock.Lock()
		if m.closed {
			return // TODO: return error
		}
		_, exists := m.ops[o.uid]
		if exists {
			m.lock.Unlock()
			continue
		}
		m.ops[o.uid] = o
		m.lock.Unlock()
		break
	}
}

func (m *operationManager) remove(o *operation) {
	m.lock.Lock()
	delete(m.ops, o.uid)
	m.lock.Unlock()
}

func (m *operationManager) shutdown() {
	m.lock.Lock()
	var wg sync.WaitGroup
	wg.Add(len(m.ops))
	m.closed = true
	for _, o := range m.ops {
		go func(op *operation) {
			// TODO: op.close()
			wg.Done()
		}(o)
	}
	m.lock.Unlock()
	wg.Wait()
}

func (m *operationManager) Range(fn func(o *operation)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, v := range m.ops {
		fn(v)
	}
}

type connection struct {
	*Server
	id            string
	ws            *websocket.Conn
	wg            sync.WaitGroup
	inbox         chan *Message
	outbox        chan *Message
	closed        bool
	keepAlive     chan time.Time
	connError     error
	connErrorLock sync.Mutex

	// operations
	operations *operationManager
}

func (s *Server) newConnection(ws *websocket.Conn) *connection {
	c := new(connection)
	c.Server = s
	c.ws = ws
	c.connections.add(c)
	c.inbox = make(chan *Message)
	c.outbox = make(chan *Message)
	c.keepAlive = make(chan time.Time, 1)
	c.operations = c.newOperationManager()
	return c
}

func (c *connection) closeWithError(err error) {
	defer func() {
		_ = recover()
	}()

	c.connErrorLock.Lock()
	if c.closed {
		c.connErrorLock.Unlock()
		return
	}
	c.closed = true
	if c.connError == nil {
		c.connError = err
	}
	c.connErrorLock.Unlock()
	c.operations.shutdown()
	close(c.outbox)
	c.wg.Wait()
}

func (c *connection) close() {
	c.closeWithError(nil)
}

func (c *connection) run() {
	c.wg.Add(1)
	go c.readLoop()

	msg, more := <-c.inbox
	if more {
		if msg.Type != MessageTypeGQLConnectionInit {
			err := errors.New("client failed to send connection_init")
			c.log.Error(err.Error())
			c.closeWithError(err)
		} else {
			c.log.Info("Received message from client: connection_init")
		}
	}

	c.wg.Add(1)
	go c.writeLoop()

	_ = c.send(&Message{
		Type: MessageTypeGQLConnectionAck,
	})

	c.wg.Add(1)
	go c.inboxLoop()

	c.wg.Add(1)
	go c.keepAliveLoop()

	c.wg.Wait()
	c.connections.remove(c)
	c.log.Info("Connection terminated")
}

func (c *connection) keepAliveLoop() {
	defer func() {
		c.wg.Done()
	}()
	if !c.cfg.EnableKeepAlive {
		return
	}
	defer c.log.Info("Keep alive routine terminated")
	c.log.Info("Keep alive routine commenced")
	for {
		select {
		case _, more := <-c.keepAlive:
			if !more {
				return
			}
		case <-time.After(time.Second * 5):
			_ = c.send(&Message{
				Type: MessageTypeGQLConnectionKeepAlive,
			})
		}
	}
}

func (c *connection) readLoop() {
	c.log.Info("Read loop commenced")
	defer func() {
		close(c.inbox)
		c.log.Info("Read loop terminated")
		c.wg.Done()
	}()
	for {
		msg := new(Message)
		err := c.ws.ReadJSON(msg)
		if err != nil {
			err = fmt.Errorf("failed to read message: %v", err)
			c.closeWithError(err)
			return
		}
		c.inbox <- msg
	}
}

func (c *connection) writeDeadline() time.Time {
	if c.cfg.WriteTimeout == 0 {
		return time.Time{}
	}
	return time.Now().Add(c.cfg.WriteTimeout)
}

func (c *connection) send(message *Message) error {
	var err error
	func() {
		defer func() {
			r := recover()
			if r != nil {
				err = c.connError
			}
		}()
		c.outbox <- message
	}()
	return err
}

func (c *connection) writeLoop() {
	c.log.Info("Write loop commenced")
	defer func() {
		c.log.Info("Write loop terminated")
		close(c.keepAlive)
		_ = c.ws.WriteMessage(websocket.CloseMessage, []byte{})
		c.ws.Close()
		c.wg.Done()
	}()
	for {
		message, more := <-c.outbox
		if !more {
			return
		}
		if c.cfg.WriteTimeout != 0 {
			err := c.ws.SetWriteDeadline(c.writeDeadline())
			if err != nil {
				c.log.Error(fmt.Sprintf("Failed to set write deadline: %v", err))
				c.closeWithError(err)
				return
			}
		}
		err := c.ws.WriteJSON(message)
		if err != nil {
			c.log.Error(fmt.Sprintf("Failed to send message: %v", err))
			c.closeWithError(err)
			return
		}

		c.keepAlive <- time.Now()
		c.log.Info(fmt.Sprintf("Sent a message to the server: %s ", message.Type))
	}
}

func (c *connection) inboxLoop() {
	c.log.Info("Inbox loop commenced")
	defer func() {
		c.log.Info("Inbox loop terminated")
		c.wg.Done()
	}()
	for {
		message, more := <-c.inbox
		if !more {
			return
		}
		switch message.Type {
		case MessageTypeGQLStart:
			c.start(message)
		case MessageTypeGQLStop:
			c.stop(message)
		case MessageTypeGQLConnectionTerminate:
			go c.terminate()
		default:
			err := fmt.Errorf("Received unsupported message type: %v", message.Type)
			c.log.Error(err.Error())
			_ = c.send(&Message{
				Type:    MessageTypeGQLConnectionError,
				Payload: err.Error(),
			})
		}
	}
}

func (c *connection) terminate() {
	c.log.Info("Terminating connection at the request of the client")
	c.close()
}
