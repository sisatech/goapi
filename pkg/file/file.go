package file

import (
	"io"
	"os"
	"time"
)

// File represents a file from the filesystem.
type File interface {

	// Name returns the base name of the file, not a
	// full path (see filepath.Base).
	Name() string

	// Size returns the size of the file in bytes. If
	// the file represents a directory the size returned
	// should be zero.
	Size() int

	// ModTime returns the time the file was most
	// recently modified.
	ModTime() time.Time

	// Read implements io.Reader to retrieve file
	// contents.
	Read(p []byte) (n int, err error)

	// Close implements io.Closer.
	Close() error

	// IsDir returns true if the File represents a
	// directory.
	IsDir() bool
}

// Open mimics the os.Open function but returns an
// implementation of File.
func Open(path string) (File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return CustomFile(CustomFileArgs{
		Name:       fi.Name(),
		Size:       int(fi.Size()),
		ModTime:    fi.ModTime(),
		IsDir:      fi.IsDir(),
		ReadCloser: f,
	}), nil
}

// CustomFileArgs takes all elements that need to be provided
// to the CustomFile function.
type CustomFileArgs struct {
	Name       string
	Size       int
	ModTime    time.Time
	IsDir      bool
	ReadCloser io.ReadCloser
}

// CustomFile makes it possible to construct a custom file
// that implements the File interface without necessarily
// being backed by an actual file on the filesystem.
func CustomFile(args CustomFileArgs) File {
	return &customFile{
		name:    args.Name,
		size:    args.Size,
		modTime: args.ModTime,
		isDir:   args.IsDir,
		rc:      args.ReadCloser,
	}
}

type customFile struct {
	name    string
	size    int
	modTime time.Time
	isDir   bool
	rc      io.ReadCloser
}

func (f *customFile) Name() string {
	return f.name
}

func (f *customFile) Size() int {
	return f.size
}

func (f *customFile) ModTime() time.Time {
	return f.modTime
}

func (f *customFile) IsDir() bool {
	return f.isDir
}

func (f *customFile) Read(p []byte) (n int, err error) {
	return f.rc.Read(p)
}

func (f *customFile) Close() error {
	if f.rc != nil {
		return f.rc.Close()
	}
	return nil
}
