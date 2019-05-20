package goapi

import (
	"encoding/json"

	"github.com/sisatech/goapi/pkg/graphqlws"
	"github.com/sisatech/goapi/pkg/objects"
)

// VMSubscription ...
func (c *Client) VMSubscription(id string, dataCallback func(payload *objects.VM, errs []graphqlws.GQLError), errCallback func(err error)) (*graphqlws.Subscription, error) {

	var dc func(payload *graphqlws.GQLDataPayload)
	dc = func(payload *graphqlws.GQLDataPayload) {
		if payload.Data == nil {
			dataCallback(nil, payload.Errors)
			return
		}

		type responseContainer struct {
			Data struct {
				VM objects.VM `json:"vm"`
			} `json:"data"`
		}

		resp := new(responseContainer)

		b, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, resp)
		if err != nil {
			panic(err)
		}

		dataCallback(&resp.Data.VM, payload.Errors)
	}

	subscription, err := c.subscriptions.Subscription(&graphqlws.SubscriptionConfig{
		Query: `
			subscription($id: ID!){
				vm(id: $id){
					args
					author
					binary
					cpus
					created
					date
					disk
					download
					env
					hostname
					id
					instance
					kernel
					logFile
					name
					networks {
					gateway
					http {
						address
						port
					}
					https {
						address
						port
					}
					ip
					mask
					name
					tcp {
						address
						port
					}
					udp {
						address
						port
					}
					}
					platform
					ram
					redirects {
					address
					source
					}
					serial {
					cursor
					data
					more
					}
					source {
					checksum
					filesystem
					icon
					job
					name
					type
					}
					stateLog
					status
					summary
					url
					version
				}
			}
		`,
		Variables: map[string]interface{}{
			"id": id,
		},
		DataCallback:  dc,
		ErrorCallback: errCallback,
	})
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// JobSubscription ...
func (c *Client) JobSubscription(id string, dataCallback func(payload *objects.Job, errs []graphqlws.GQLError), errCallback func(err error)) (*graphqlws.Subscription, error) {

	var dc func(payload *graphqlws.GQLDataPayload)
	dc = func(payload *graphqlws.GQLDataPayload) {
		if payload.Data == nil {
			dataCallback(nil, payload.Errors)
			return
		}

		type responseContainer struct {
			Data struct {
				Job objects.Job `json:"job"`
			} `json:"data"`
		}

		resp := new(responseContainer)

		b, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, resp)
		if err != nil {
			panic(err)
		}

		dataCallback(&resp.Data.Job, payload.Errors)
	}

	subscription, err := c.subscriptions.Subscription(&graphqlws.SubscriptionConfig{
		Query: `
			subscription($id: String!){
				job (id: $id) {
					description
					id
					logFilePath
					logPlainFilePath
					name
					progress {
						error
						finished
						progress
						started
						status
						total
						units
					}
				}
			}	
		`,
		Variables: map[string]interface{}{
			"id": id,
		},
		DataCallback:  dc,
		ErrorCallback: errCallback,
	})
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// DefaultsSubscription ..
func (c *Client) DefaultsSubscription(dataCallback func(payload *objects.Defaults, errs []graphqlws.GQLError),
	errCallback func(err error)) (*graphqlws.Subscription, error) {

	var dc func(payload *graphqlws.GQLDataPayload)
	dc = func(payload *graphqlws.GQLDataPayload) {

		if payload.Data != nil {
			dataCallback(nil, payload.Errors)
		}

		type responseContainer struct {
			Data struct {
				Defaults objects.Defaults `json:"defaults"`
			} `json:"data"`
		}

		resp := new(responseContainer)

		b, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, resp)
		if err != nil {
			panic(err)
		}

		dataCallback(&resp.Data.Defaults, payload.Errors)
	}

	subscription, err := c.subscriptions.Subscription(&graphqlws.SubscriptionConfig{
		Query: `
                        subscription {
                                defaults {
                                        kernel
                                        platform
                                }
                        }`,
		DataCallback:  dc,
		ErrorCallback: errCallback,
	})
	if err != nil {
		return nil, err
	}

	return subscription, nil
}
