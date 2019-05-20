package goapi

import (
	"fmt"
	"time"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

// App ..
type App struct {
	bucket *Bucket
	name   string
}

// AppList ..
type AppList struct {
	PageInfo objects.PageInfo
	Items    []AppListItem
}

// AppListItem ..
type AppListItem struct {
	Cursor string
	App    App
}

// Name ..
func (a *App) Name() string {
	return a.name
}

// Version ..
func (a *App) Version(ref string) (*Version, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                app(name: "%s") {
                                        version(ref: "%s") {
                                                id
                                                uploadedTimeplate
                                        }
                                }
                        }
                }
        `, a.bucket.Name(), a.Name(), ref))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := a.bucket.g.client.Run(a.bucket.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &Version{
		app:               a,
		id:                resp.Bucket.App.Version.ID,
		uploadedTimeplate: time.Unix(int64(resp.Bucket.App.Version.UploadedTimeplate), 0),
	}, nil
}

// VersionList ..
func (a *App) VersionList(curs *Cursor) (*VersionList, error) {

	var vd, v string
	if curs != nil {
		vd, v = curs.Strings()
	}

	req := graphql.NewRequest(fmt.Sprintf(`
                query%s {
                        bucket(name: "%s") {
                                app(name: "%s") {
                                        versionsList%s {
                                                edges {
                                                        cursor
                                                        node {
                                                                id
                                                                uploadedTimeplate
                                                        }
                                                }
                                                pageInfo {
                                                        hasNextPage
                                                        hasPreviousPage
                                                        endCursor
                                                        startCursor
                                                }
                                        }
                                }
                        }
                }
        `, vd, a.bucket.Name(), a.Name(), v))
	if curs != nil {
		curs.AddToRequest(req)
	}

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := a.bucket.g.client.Run(a.bucket.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	out := new(VersionList)
	out.PageInfo = resp.Bucket.App.VersionsList.PageInfo
	out.Items = make([]VersionListItem, 0)

	for _, v := range resp.Bucket.App.VersionsList.Edges {
		out.Items = append(out.Items, VersionListItem{
			Cursor: v.Cursor,
			Version: Version{
				app:               a,
				id:                v.Node.ID,
				uploadedTimeplate: time.Unix(int64(v.Node.UploadedTimeplate), 0),
			},
		})
	}

	return out, nil
}

// Latest ..
func (a *App) Latest() (*Version, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                app(name: "%s") {
                                        latest {
                                                id
                                                uploadedTimeplate
                                        }
                                }
                        }
                }
        `, a.bucket.Name(), a.Name()))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := a.bucket.g.client.Run(a.bucket.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &Version{
		app:               a,
		id:                resp.Bucket.App.Latest.ID,
		uploadedTimeplate: time.Unix(int64(resp.Bucket.App.Latest.UploadedTimeplate), 0),
	}, nil
}

// Authorization ..
func (a *App) Authorization() (*objects.Authorization, error) {

	req := graphql.NewRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                app(name: "%s") {
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
                }
        `, a.bucket.Name(), a.Name()))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := a.bucket.g.client.Run(a.bucket.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Bucket.App.Authorization, nil
}
