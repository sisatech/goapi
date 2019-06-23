package goapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/graphqlws"
	"github.com/sisatech/goapi/pkg/objects"
)

// Repository provides access to APIs specific to a single repository.
type Repository struct {
	mgr      *RepositoriesManager
	name     string
	host     string
	subsDoer *graphqlws.Client
	hdr      http.Header
}

func (r *Repository) init() error {

	m := make(map[string][]string)
	m["Vorteil"] = []string{r.name}

	r.hdr = http.Header(m)

	// var err error
	// r.subsDoer, err = graphqlws.NewClient(r.mgr.c.ctx, &graphqlws.ClientConfig{
	// 	Address: r.mgr.c.cfg.Address,
	// 	Path:    "subscriptions",
	// 	Header:  r.hdr,
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (r *Repository) newRequest(str string) *graphql.Request {
	req := graphql.NewRequest(str)
	req.Header = r.hdr
	return req
}

// // Delete an bucket within the repository.
// func (r *Repository) DeleteBucket(name string) error {
//
// 	req := r.newRequest(fmt.Sprintf(`
// 		mutation {
// 			deleteBucket(name: "%s")
// 		}
// 	`, name))
//
// 	type responseContainer struct {
// 		DeleteBucket bool `json:"deleteBucket"`
// 	}
// 	resp := new(responseContainer)
// 	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// // Delete an app within the repository.
// func (r *Repository) DeleteApp(bucket, app string) error {
//
// 	req := r.newRequest(fmt.Sprintf(`
// 		mutation {
// 			deleteApp(bucketName: "%s", appName: "%s")
// 		}
// 	`, bucket, app))
//
// 	type responseContainer struct {
// 		DeleteApp bool `json:"deleteApp"`
// 	}
// 	resp := new(responseContainer)
// 	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// // DeleteAppVersion deletes a specific app version from within a bucket.
// func (r *Repository) DeleteVersion(bucket, app, version string) error {
// 	req := r.newRequest(fmt.Sprintf(`
// 		mutation {
// 			deleteAppVersion(bucketName: "%s", appName: "%s", reference: "%s")
// 		}
// 	`, bucket, app, version))
//
// 	type responseContainer struct {
// 		DeleteAppVersion bool `json:"deleteAppVersion"`
// 	}
// 	resp := new(responseContainer)
// 	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// Download an object (app/version) from the repository.
func (r *Repository) Download(bucket, app, version string, w io.Writer) error {

	var versionString = fmt.Sprintf(`version(ref: "%s")`, version)
	if version == "" {
		versionString = fmt.Sprintf(`latest`)
	}

	req := r.newRequest(fmt.Sprintf(`
		query {
			bucket(name: "%s") {
				app(name: "%s") {
					%s {
						file {
							downloadURL
						}
					}
				}
			}
		}
	`, bucket, app, versionString))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}
	resp := new(responseContainer)
	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
	if err != nil {
		return err
	}

	var downloadURL string
	if version == "" {
		downloadURL = resp.Bucket.App.Latest.File.DownloadURL
	} else {
		downloadURL = resp.Bucket.App.Version.File.DownloadURL
	}

	re, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", r.host, downloadURL), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(re)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status code: %v - %s",
			res.StatusCode, res.Status)
	}

	_, err = io.Copy(w, res.Body)
	if err != nil {
		return err
	}

	return nil
}

//  Upload an object (app/version) to the repository.
func (r *Repository) Upload() error {
	return nil
}

// NewBucket creates a new bucket within the repository.
func (r *Repository) NewBucket(name string) error {

	req := r.newRequest(fmt.Sprintf(`
		mutation {
			newBucket(name: "%s") {
				name
			}
		}
	`, name))

	type responseContainer struct {
		NewBucket objects.Bucket `json:"newBucket"`
	}
	resp := new(responseContainer)
	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
	if err != nil {
		return err
	}

	return nil
}

// GetBucket ..
func (r *Repository) GetBucket(name string) (*Bucket, error) {

	req := r.newRequest(fmt.Sprintf(`
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
	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		r:    r,
		name: resp.Bucket.Name,
	}, nil
}

// ListBuckets returns a list of buckets within the repository.
func (r *Repository) ListBuckets(cursor *Cursor) (*BucketList, error) {

	var vd, v string
	if cursor != nil {
		vd, v = cursor.Strings()
	}

	req := r.newRequest(fmt.Sprintf(`
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
	if cursor != nil {
		cursor.AddToRequest(req)
	}

	type responseContainer struct {
		ListBuckets objects.BucketsConnection `json:"listBuckets"`
	}
	resp := new(responseContainer)
	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	// out := make([]string, 0)
	// for _, b := range resp.ListBuckets.Edges {
	// 	out = append(out, b.Node.Name)
	// }

	out := new(BucketList)
	out.PageInfo = resp.ListBuckets.PageInfo
	out.Items = make([]BucketListItem, 0)

	for _, b := range resp.ListBuckets.Edges {
		out.Items = append(out.Items, BucketListItem{
			Cursor: b.Cursor,
			Bucket: Bucket{
				r:    r,
				name: b.Node.Name,
			},
		})
	}

	return out, nil
}

// // ListApps returns a list of apps within a bucket.
// func (r *Repository) ListApps(bucket string, cursor *Cursor) ([]string, error) {
//
// 	var vd, v string
// 	if cursor != nil {
// 		vd, v = cursor.Strings()
// 	}
//
// 	req := r.newRequest(fmt.Sprintf(`
// 		query%s {
// 			bucket(name: "%s") {
// 				appsList%s {
// 					edges {
// 						node {
// 							name
// 						}
// 					}
// 				}
// 			}
// 		}
// 	`, vd, bucket, v))
//
// 	if cursor != nil {
// 		cursor.AddToRequest(req)
// 	}
//
// 	type responseContainer struct {
// 		Bucket objects.Bucket `json:"bucket"`
// 	}
// 	resp := new(responseContainer)
// 	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	out := make([]string, 0)
// 	for _, a := range resp.Bucket.AppsList.Edges {
// 		out = append(out, a.Node.Name)
// 	}
//
// 	return out, nil
// }

// VersionSummary contains basic info about an app version within an repository.
// type VersionSummary struct {
// 	ID       string
// 	Tag      string
// 	Uploaded time.Time
// }

// // ListVersions returns a list of versions of an app.
// func (r *Repository) ListVersions(bucket, app string, cursor *Cursor) ([]VersionSummary, error) {
//
// 	var vd, v string
// 	if cursor != nil {
// 		vd, v = cursor.Strings()
// 	}
//
// 	req := r.newRequest(fmt.Sprintf(`
// 		query%s {
// 			bucket(name: "%s") {
// 				app(name: "%s") {
// 					versionsList%s {
// 						edges {
// 							node {
// 								id
// 								tag
// 								uploadedTimeplate
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	`, vd, bucket, app, v))
//
// 	type responseContainer struct {
// 		Bucket objects.Bucket `json:"bucket"`
// 	}
// 	resp := new(responseContainer)
// 	err := r.mgr.c.graphql.Run(r.mgr.c.ctx, req, &resp)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	out := make([]VersionSummary, 0)
// 	for _, v := range resp.Bucket.App.VersionsList.Edges {
// 		out = append(out, VersionSummary{
// 			ID:       v.Node.ID,
// 			Tag:      v.Node.Tag,
// 			Uploaded: time.Unix(int64(v.Node.UploadedTimeplate), 0),
// 		})
// 	}
//
// 	return out, nil
// }
