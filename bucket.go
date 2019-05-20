package goapi

import (
	"encoding/json"
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/graphqlws"
	"github.com/sisatech/goapi/pkg/objects"
)

// Bucket ..
type Bucket struct {
	g    *Client
	name string
}

// BucketList ..
type BucketList struct {
	PageInfo objects.PageInfo
	Items    []BucketListItem
}

// BucketListItem ..
type BucketListItem struct {
	Cursor string
	Bucket Bucket
}

// Name ..
func (b *Bucket) Name() string {
	return b.name
}

// App ..
func (b *Bucket) App(name string) (*App, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                app(name: "%s") {
                                        name
                                }
                        }
                }
        `, b.Name(), name))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := b.g.client.Run(b.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &App{
		bucket: b,
		name:   resp.Bucket.App.Name,
	}, nil
}

// AppList ..
func (b *Bucket) AppList(curs *Cursor) (*AppList, error) {

	var vd, v string
	if curs != nil {
		vd, v = curs.Strings()
	}

	req := graphql.NewRequest(fmt.Sprintf(`
                query%s {
                        bucket(name: "%s") {
                                appsList%s {
                                        edges {
                                                cursor
                                                node
                                        }
                                        pageInfo {
                                                endCursor
                                                startCursor
                                                hasNextPage
                                                hasPreviousPage
                                        }
                                }
                        }
                }
        `, vd, b.Name(), v))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := b.g.client.Run(b.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	out := new(AppList)
	out.PageInfo = resp.Bucket.AppsList.PageInfo
	out.Items = make([]AppListItem, 0)

	for _, a := range resp.Bucket.AppsList.Edges {
		out.Items = append(out.Items, AppListItem{
			Cursor: a.Cursor,
			App: App{
				bucket: b,
				name:   a.Node.Name,
			},
		})
	}

	return out, nil
}

// Authorization ..
func (b *Bucket) Authorization() (*objects.Authorization, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                authorization {
                                        id
                                        owner
                                        acls {
                                                group
                                                action
                                        }
                                }
                        }
                }
        `, b.Name()))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := b.g.client.Run(b.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Bucket.Authorization, nil
}

// Icon ..
func (b *Bucket) Icon() (*objects.Fragment, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                icon {
                                        id
                                        md5
                                        downloadURL
                                        uploadURL
                                }
                        }
                }
        `, b.Name()))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := b.g.client.Run(b.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Bucket.Icon, nil
}

// GetBucket ..
func (c *Client) GetBucket(name string) (*Bucket, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                name
                        }
                }
        `, name))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		g:    c,
		name: resp.Bucket.Name,
	}, nil
}

// ListBuckets ..
func (c *Client) ListBuckets(curs *Cursor) (*BucketList, error) {

	var vd, v string
	if curs != nil {
		vd, v = curs.Strings()
	}

	req := graphql.NewRequest(fmt.Sprintf(`
                query%s {
                        listBuckets%s {
                                edges {
                                        cursor
                                        node {
                                                name
                                        }
                                }
                                pageInfo {
                                        endCursor
                                        startCursor
                                        hasNextPage
                                        hasPreviousPage
                                }
                        }
                }
        `, vd, v))
	if curs != nil {
		curs.AddToRequest(req)
	}

	type responseContainer struct {
		ListBucket objects.BucketsConnection `json:"listBuckets"`
	}

	resp := new(responseContainer)
	err := c.client.Run(c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	out := new(BucketList)
	out.PageInfo = resp.ListBucket.PageInfo
	out.Items = make([]BucketListItem, 0)

	for _, b := range resp.ListBucket.Edges {
		out.Items = append(out.Items, BucketListItem{
			Cursor: b.Cursor,
			Bucket: Bucket{
				g:    c,
				name: b.Node.Name,
			},
		})
	}

	return out, nil
}

// ListBucketsSubscription ..
func (c *Client) ListBucketsSubscription(dataCallback func([]string, []graphqlws.GQLError),
	errCallback func(error)) (*graphqlws.Subscription, error) {

	dc := func(payload *graphqlws.GQLDataPayload) {

		if payload.Data == nil {
			dataCallback(nil, payload.Errors)
			return
		}

		type responseContainer struct {
			Data struct {
				ListBuckets objects.BucketsConnection `json:"listBuckets"`
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

		out := make([]string, 0)
		for _, b := range resp.Data.ListBuckets.Edges {
			out = append(out, b.Node.Name)
		}

		dataCallback(out, payload.Errors)
	}

	subscription, err := c.subscriptions.Subscription(&graphqlws.SubscriptionConfig{
		Query: `
			subscription {
				listBuckets {
					edges {
						node {
							name
						}
					}
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
