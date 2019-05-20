package file

import (
	"os"
	"time"
)

// finfo exists to implement os.FileInfo for this package's
// FileInfo function.
type finfo struct {
	name    string
	size    int64
	modtime time.Time
	mode    os.FileMode
}

func (fi *finfo) Name() string {
	return fi.name
}

func (fi *finfo) Size() int64 {
	return fi.size
}

func (fi *finfo) ModTime() time.Time {
	return fi.modtime
}

func (fi *finfo) Mode() os.FileMode {
	return fi.mode
}

func (fi *finfo) IsDir() bool {
	return fi.mode.IsDir()
}

func (fi *finfo) Sys() interface{} {
	return nil
}

// Info produces an implementation of os.FileInfo from a an
// implementation of File.
func Info(f File) os.FileInfo {
	mode := os.ModePerm
	if f.IsDir() {
		mode |= os.ModeDir
	}
	return &finfo{
		name:    f.Name(),
		size:    int64(f.Size()),
		modtime: f.ModTime(),
		mode:    mode,
	}
}
