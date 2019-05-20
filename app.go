package goapi

import "github.com/sisatech/goapi/pkg/objects"

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

// Latest ..
func (a *App) Version() (*Version, error) {
	return nil, nil
}
