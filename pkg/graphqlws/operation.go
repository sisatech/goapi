package graphqlws

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

type operation struct {
	*connection
	uid           string
	id            string
	query         string
	operationName string
	variables     map[string]interface{}

	// updater
	lastMD5    string
	lastUpdate time.Time
	updated    chan time.Time
	update     Update
}

func (c *connection) start(m *Message) {
	c.log.Info("Client commencing a new operation")

	gqlError := func(err error) {
		_ = c.send(&Message{
			Type:    MessageTypeGQLError,
			ID:      m.ID,
			Payload: err.Error(),
		})
		_ = c.send(&Message{
			Type: MessageTypeGQLComplete,
			ID:   m.ID,
		})
	}

	o := new(operation)
	o.connection = c
	o.updated = make(chan time.Time)
	o.id = m.ID

	// process arguments
	if m.Payload == nil {
		gqlError(errors.New("cannot accept nil payload on GQL_START message"))
		return
	}
	p, ok := m.Payload.(map[string]interface{})
	if !ok {
		gqlError(errors.New("payload is invalid for a GQL_START message"))
		return
	}
	for k, v := range p {
		switch k {
		case "operationName":
			o.operationName, ok = v.(string)
			if !ok {
				gqlError(errors.New("payload field 'operationName' has invalid type"))
				return
			}
		case "query":
			o.query, ok = v.(string)
			if !ok {
				gqlError(errors.New("payload field 'query' has invalid type"))
				return
			}
		case "variables":
			o.variables, ok = v.(map[string]interface{})
			if !ok {
				gqlError(errors.New("payload field 'variables' has invalid type"))
				return
			}
		case "extensions":
			extensions, ok := v.(map[string]interface{})
			if !ok {
				gqlError(errors.New("payload field 'extensions' has invalid type"))
				return
			}
			if len(extensions) > 0 {
				gqlError(errors.New("extensions are not supported"))
				return
			}
		default:
			gqlError(fmt.Errorf("payload field '%s' is unsupported", k))
			return
		}
	}

	// validate query
	document, err := parser.Parse(parser.ParseParams{
		Source: o.query,
	})
	if err != nil {
		gqlError(err)
		return
	}

	validation := graphql.ValidateDocument(&c.cfg.Schema, document, nil)
	if !validation.IsValid {
		gqlError(validation.Errors[0])
		return
	}

	// trigger first response
	o.execute()

	// initialize operation fields
	var fields []string
	var recurseForFields func(parent string, ss *ast.SelectionSet)
	recurseForFields = func(parent string, ss *ast.SelectionSet) {
		for _, selection := range ss.Selections {
			field, ok := selection.(*ast.Field)
			if !ok {
				fmt.Printf("FAILED: %v", reflect.TypeOf(selection))
				return
			}
			key := parent + field.Name.Value
			if field.SelectionSet == nil {
				fields = append(fields, key)
				return
			}
			key += "."
			recurseForFields(key, field.SelectionSet)
		}
	}

	var requiresSubscription bool
	for _, definition := range document.Definitions {
		if definition.GetKind() == "OperationDefinition" {
			od := definition.(*ast.OperationDefinition)
			if od.Operation != "subscription" {
				continue
			}
			requiresSubscription = true
			ss := od.SelectionSet
			if ss == nil {
				continue
			}
			recurseForFields("", ss)
		}
	}

	o.update = o.schema.NewUpdate(fields...)

	// terminate operation if no subscriptions are involved
	if !requiresSubscription {
		o.finish()
	} else {
		// add operation to connection
		c.operations.add(o)
		o.subscribe()
	}

}

func (c *connection) stop(m *Message) {
	c.log.Info("Client operation terminated")
	o, ok := c.operations.get(m.ID)
	if ok {
		o.complete()
	}
}

func (o *operation) finish() {
	_ = o.send(&Message{
		Type: MessageTypeGQLComplete,
		ID:   o.id,
	})
}

func (o *operation) complete() {
	o.finish()
	o.operations.remove(o)
}

func (o *operation) execute() {
	result := graphql.Do(graphql.Params{
		OperationName:  o.operationName,
		RequestString:  o.query,
		VariableValues: o.variables,
		Schema:         o.cfg.Schema,
		// TODO: Context: nil,
	})
	msg := &Message{
		Type: MessageTypeGQLData,
		ID:   o.id,
		Payload: map[string]interface{}{
			"data":   result.Data,
			"errors": result.Errors,
		},
	}
	data, _ := json.Marshal(msg)
	hash := md5.Sum(data)
	hashString := hex.EncodeToString(hash[:])
	if hashString == o.lastMD5 {
		go func(t time.Time) {
			o.updated <- t
		}(time.Now())
		return
	} else {
		o.lastMD5 = hashString
	}

	_ = o.send(msg)

	go func(t time.Time) {
		o.updated <- t
	}(time.Now())

}

func (o *operation) subscribe() {
	go o.poller()
}

func (o *operation) poller() {
	interval := o.cfg.PollingInterval
	if interval == time.Duration(0) {
		return
	}
	o.lastUpdate = time.Now()
	for {
		dt := interval - time.Since(o.lastUpdate)
		if dt < 0 {
			dt = 0
		}

		select {
		case <-time.After(dt):
			o.execute()
			time.Sleep(interval / 2)
		case t := <-o.updated:
			o.lastUpdate = t
		}
	}
}
