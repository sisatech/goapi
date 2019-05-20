package goapi

import (
	"encoding/json"

	"github.com/sisatech/goapi/pkg/graphqlws"
	"github.com/sisatech/goapi/pkg/objects"
)

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
