package goapi

import (
	"encoding/json"

	"github.com/sisatech/goapi/pkg/graphqlws"
	"github.com/sisatech/goapi/pkg/objects"
)

// ListNodesSubscription ..
func (c *Client) ListNodesSubscription(dataCallback func([]objects.Node, []graphqlws.GQLError),
	errCallback func(error)) (*graphqlws.Subscription, error) {

	dc := func(payload *graphqlws.GQLDataPayload) {

		if payload.Data == nil {
			dataCallback(nil, payload.Errors)
			return
		}

		type responseContainer struct {
			Data struct {
				ListNodes []objects.Node `json:"listNodes"`
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

		dataCallback(resp.Data.ListNodes, payload.Errors)
	}

	subscription, err := c.subscriptions.Subscription(&graphqlws.SubscriptionConfig{
		DataCallback:  dc,
		ErrorCallback: errCallback,
		Query: `
		subscription {
			listNodes {
				host
				name
				type
			}
		}`,
	})
	if err != nil {
		return nil, err
	}

	return subscription, nil
}
