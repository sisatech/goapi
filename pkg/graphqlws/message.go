package graphqlws

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MessageType is a string representing a specific type of message.
type MessageType string

// GQL message types
const (
	MessageTypeGQLConnectionInit      MessageType = "connection_init"
	MessageTypeGQLStart               MessageType = "start"
	MessageTypeGQLStop                MessageType = "stop"
	MessageTypeGQLConnectionTerminate MessageType = "connection_terminate"
	MessageTypeGQLConnectionError     MessageType = "connection_error"
	MessageTypeGQLConnectionAck       MessageType = "connection_ack"
	MessageTypeGQLData                MessageType = "data"
	MessageTypeGQLError               MessageType = "error"
	MessageTypeGQLComplete            MessageType = "complete"
	MessageTypeGQLConnectionKeepAlive MessageType = "ka"
)

// Message is the generalized form for a GraphQL over websocket message.
type Message struct {
	Payload interface{} `json:"payload,omitempty"`
	ID      string      `json:"id,omitempty"`
	Type    MessageType `json:"type"`
}

func (m *Message) GQLDataPayload() (*GQLDataPayload, error) {
	if m.Type != MessageTypeGQLData {
		return nil, fmt.Errorf("message isn't of type GQL_DATA (have %s)", m.Type)
	}
	pl := new(GQLDataPayload)
	if m.Payload == nil {
		return pl, nil
	}
	data, err := json.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	err = dec.Decode(pl)
	if err != nil {
		return nil, err
	}
	return pl, nil
}

// GQLErrorLocation contains a reference to a location in the operation's
// query string related to an error that has occurred.
type GQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// GQLError contains an error description with information meant to be useful
// for a developer that needs to debug it.
type GQLError struct {
	Locations  []GQLErrorLocation     `json:"locations,omitempty"`
	Message    string                 `json:"message"`
	Path       []string               `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func (e *GQLError) Error() string {
	s := e.Message
	// TODO: add more information from e
	return s
}

// GQLDataPayload is a generic structure for a response to a valid GraphQL
// operation.
type GQLDataPayload struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []GQLError  `json:"errors,omitempty"`
}

// GQLStartPayload is the structure for a GQL_START message payload.
type GQLStartPayload struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}
