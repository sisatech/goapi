package goapi

import (
	"fmt"
	"time"

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

	req := v.app.bucket.g.NewRequest(fmt.Sprintf(`
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
	err := v.app.bucket.g.client.Run(v.app.bucket.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Bucket.App.Version.File, nil
}

// Icon ..
func (v *Version) Icon() (*objects.PackageFragment, error) {

	req := v.app.bucket.g.NewRequest(fmt.Sprintf(`
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
	err := v.app.bucket.g.client.Run(v.app.bucket.g.ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Bucket.App.Version.Icon, nil
}

// Tag ..
func (v *Version) Tag() (string, error) {

	req := v.app.bucket.g.NewRequest(fmt.Sprintf(`
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
	err := v.app.bucket.g.client.Run(v.app.bucket.g.ctx, req, &resp)
	if err != nil {
		return "", err
	}

	return resp.Bucket.App.Version.Tag, nil
}
