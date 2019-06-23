package goapi

import (
	"fmt"
	"time"

	"github.com/machinebox/graphql"
	"github.com/sisatech/goapi/pkg/objects"
)

// Version ..
type Version struct {
	app               *App
	id                string
	uploadedTimeplate time.Time
}

// VersionList ..
type VersionList struct {
	PageInfo objects.PageInfo
	Items    []VersionListItem
}

// VersionListItem ..
type VersionListItem struct {
	Cursor  string
	Version Version
}

// ID ..
func (v *Version) ID() string {
	return v.id
}

// UploadedTime ..
func (v *Version) UploadedTime() time.Time {
	return v.uploadedTimeplate
}

// File ..
func (v *Version) File() (*objects.PackageFragment, error) {

	req := v.app.bucket.r.newRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                app(name: "%s") {
                                        version(ref: "%s") {
                                                file {
                                                        md5
                                                        url
                                                        size
                                                        downloadURL
                                                }
                                        }
                                }
                        }
                }
        `, v.app.bucket.Name(), v.app.Name(), v.ID()))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := v.app.bucket.r.mgr.c.graphql.Run(v.app.bucket.r.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Bucket.App.Version.File, nil
}

// Icon ..
func (v *Version) Icon() (*objects.PackageFragment, error) {

	req := v.app.bucket.r.newRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                app(name: "%s") {
                                        version(ref: "%s") {
                                                icon {
                                                        md5
                                                        url
                                                        size
                                                        downloadURL
                                                }
                                        }
                                }
                        }
                }
        `, v.app.bucket.Name(), v.app.Name(), v.ID()))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := v.app.bucket.r.mgr.c.graphql.Run(v.app.bucket.r.mgr.c.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Bucket.App.Version.Icon, nil
}

// Tag ..
func (v *Version) Tag() (string, error) {

	req := v.app.bucket.r.newRequest(fmt.Sprintf(`
                query {
                        bucket(name: "%s") {
                                app(name: "%s") {
                                        version(ref: "%s") {
                                                tag
                                        }
                                }
                        }
                }
        `, v.app.bucket.Name(), v.app.Name(), v.ID()))

	type responseContainer struct {
		Bucket objects.Bucket `json:"bucket"`
	}

	resp := new(responseContainer)
	err := v.app.bucket.r.mgr.c.graphql.Run(v.app.bucket.r.mgr.c.ctx, req, &resp)
	if err != nil {
		return "", err
	}

	return resp.Bucket.App.Version.Tag, nil
}

// SetTag ..
func (v *Version) SetTag(tag string) error {

	req := graphql.NewRequest(fmt.Sprintf(`
		mutation {
			tagApp(bucketName: "%s", appName: "%s", reference: "%s", tag: "%s")
		}
	`, v.app.bucket.Name(), v.app.Name(), v.ID(), tag))

	type responseContainer struct {
		TagApp string `json:"tagApp"`
	}

	resp := new(responseContainer)
	err := v.app.bucket.r.mgr.c.graphql.Run(v.app.bucket.r.mgr.c.ctx, req, &resp)
	if err != nil {
		return err
	}

	return nil
}

// RemoveTag ..
func (v *Version) RemoveTag() error {

	req := graphql.NewRequest(fmt.Sprintf(`
		mutation {
			tagApp(bucketName: "%s", appName: "%s", reference: "%s", tag: "")
		}
	`, v.app.bucket.Name(), v.app.Name(), v.ID()))

	type responseContainer struct {
		TagApp string `json:"tagApp"`
	}

	resp := new(responseContainer)
	err := v.app.bucket.r.mgr.c.graphql.Run(v.app.bucket.r.mgr.c.ctx, req, &resp)
	if err != nil {
		return err
	}

	return nil
}

// Delete ..
func (v *Version) Delete() error {

	req := v.app.bucket.r.newRequest(fmt.Sprintf(`
		mutation {
			deleteAppVersion(bucketName: "%s", appName: "%s", reference: "%s") {
				name
			}
		}
	`, v.app.bucket.Name(), v.app.Name(), v.ID()))

	type responseContainer struct {
		DeleteAppVersion objects.App `json:"deleteAppVersion"`
	}

	resp := new(responseContainer)
	err := v.app.bucket.r.mgr.c.graphql.Run(v.app.bucket.r.mgr.c.ctx, req, &resp)
	return err
}

// Germ returns a string that can be used to identify this app/version. This can
// be used in operations such as build, or run.
func (v *Version) Germ() string {
	return fmt.Sprintf("%s:%s/%s/%s", v.app.bucket.r.name, v.app.bucket.Name(),
		v.app.Name(), v.ID())
}
