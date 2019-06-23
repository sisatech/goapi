package goapi

import (
	"fmt"
	"time"

	"github.com/sisatech/goapi/pkg/objects"
)

// App represents an application within a bucket in a Vorteil repository.
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

// Name returns the protected 'name' field of the App.
func (a *App) Name() string {
	return a.name
}

// Version fetches a specific version of an App. The 'ref' argument can be
// either the hash or tag of the desired version. If the specified version could
// not be found, or the user has insufficient permissions, an error will be
// returned. Exactly one returned argument will be non-nil.
func (a *App) Version(ref string) (*Version, error) {

	req := a.bucket.r.newRequest(fmt.Sprintf(`
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
	err := a.bucket.r.mgr.c.graphql.Run(a.bucket.r.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &Version{
		app:               a,
		id:                resp.Bucket.App.Version.ID,
		uploadedTimeplate: time.Unix(int64(resp.Bucket.App.Version.UploadedTimeplate), 0),
	}, nil
}

// VersionList fetches a list of all versions of the App. The 'curs' optional
// argument allows for pagination information to be passed to the request.
// Exactly one returned argument will be non-nil.
func (a *App) VersionList(curs *Cursor) (*VersionList, error) {

	var vd, v string
	if curs != nil {
		vd, v = curs.Strings()
	}

	req := a.bucket.r.newRequest(fmt.Sprintf(`
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
	err := a.bucket.r.mgr.c.graphql.Run(a.bucket.r.mgr.c.ctx, req, &resp)
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

// Latest fetches the most recently uploaded version of an App.
// Exactly one returned argument will be non-nil.
func (a *App) Latest() (*Version, error) {

	req := a.bucket.r.newRequest(fmt.Sprintf(`
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
	err := a.bucket.r.mgr.c.graphql.Run(a.bucket.r.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &Version{
		app:               a,
		id:                resp.Bucket.App.Latest.ID,
		uploadedTimeplate: time.Unix(int64(resp.Bucket.App.Latest.UploadedTimeplate), 0),
	}, nil
}

// Authorization returns a list of all Access Control Rules defined for the App,
// where the caller has adequate permissions to do so. Exactly one returned
// argument will be non-nil.
func (a *App) Authorization() (*objects.Authorization, error) {

	req := a.bucket.r.newRequest(fmt.Sprintf(`
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
	err := a.bucket.r.mgr.c.graphql.Run(a.bucket.r.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Bucket.App.Authorization, nil
}

// Delete deletes the App from the repository.
func (a *App) Delete() error {

	req := a.bucket.r.newRequest(fmt.Sprintf(`
		mutation {
			deleteApp(bucketName: "%s", appName: "%s")
		}
	`, a.bucket.Name(), a.Name()))

	type responseContainer struct {
		DeleteApp bool `json:"deleteApp"`
	}

	resp := new(responseContainer)
	err := a.bucket.r.mgr.c.graphql.Run(a.bucket.r.mgr.c.ctx, req, &resp)
	return err
}

// Germ returns a string that can be used to identify this app. This can
// be used in operations such as build, or run.
func (a *App) Germ() string {
	return fmt.Sprintf("%s:%s/%s", a.bucket.r.name, a.bucket.Name(), a.Name())
}
