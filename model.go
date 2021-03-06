package parcel

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"time"
)

//go:generate counterfeiter -fake-name FileSystem -o ./fake/FileSystem.go . FileSystem
//go:generate counterfeiter -fake-name File -o ./fake/File.go . File
//go:generate counterfeiter -fake-name Composer -o ./fake/Composer.go . Composer
//go:generate counterfeiter -fake-name Compressor -o ./fake/Compressor.go . Compressor

// FileSystem provides primitives to work with the underlying file system
type FileSystem interface {
	// Walk walks the file tree rooted at root, calling walkFn for each file or
	// directory in the tree, including root.
	Walk(dir string, fn filepath.WalkFunc) error
	// OpenFile is the generalized open call; most users will use Open
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
}

// Composer composes the resources
type Composer interface {
	// Compose composes from an archive
	Compose(bundle *Bundle) error
}

// Compressor compresses given resource
type Compressor interface {
	// Compress compresses given source
	Compress(fileSystem FileSystem) (*Bundle, error)
}

// Bundle represents a bundled resource
type Bundle struct {
	// Name of the resource
	Name string
	// Length returns the count of files in the bundle
	Length int
	// Body of the resource
	Body []byte
}

// File is the bundle file
type File interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

// Binary represents a resource blob content
type Binary = []byte

// Node represents a hierarchy node in the resource manager
type Node struct {
	name     string
	dir      bool
	content  Binary
	children []*Node
}

// NewNodeDir creates a new directory node
func NewNodeDir(name string) *Node {
	return &Node{
		name: name,
		dir:  true,
	}
}

// NewNodeFile creates a new file node
func NewNodeFile(name string, content Binary) *Node {
	return &Node{
		name:    name,
		dir:     false,
		content: content,
	}
}

// Name returns the base name of the file
func (n *Node) Name() string {
	return n.name
}

// Size returns the length in bytes for regular files
func (n *Node) Size() int64 {
	return int64(len(n.content))
}

// Mode returns the file mode bits
func (n *Node) Mode() os.FileMode {
	return 0
}

// ModTime returns the modification time
func (n *Node) ModTime() time.Time {
	return time.Now()
}

// IsDir returns true if the node is directory
func (n *Node) IsDir() bool {
	return n.dir
}

// Sys returns the underlying data source
func (n *Node) Sys() interface{} {
	return nil
}

// Buffer represents a *bytes.Buffer that can be closed
type Buffer struct {
	buffer *bytes.Buffer
}

// NewBuffer creates a new Buffer
func NewBuffer() *Buffer {
	return &Buffer{
		buffer: &bytes.Buffer{},
	}
}

// NewBufferWith creates a new Buffer with data
func NewBufferWith(data []byte) *Buffer {
	return &Buffer{
		buffer: bytes.NewBuffer(data),
	}
}

// Read reads the next len(p) bytes from the buffer or until the buffer is drainged
func (b *Buffer) Read(p []byte) (int, error) {
	return b.buffer.Read(p)
}

// Write appends the contents of p to the buffer, growing the buffer as needed.
func (b *Buffer) Write(p []byte) (int, error) {
	return b.buffer.Write(p)
}

// Close closes the buffer (noop).
func (b *Buffer) Close() error {
	return nil
}

// String returns the contents of the unread portion of the buffer
func (b *Buffer) String() string {
	return b.buffer.String()
}

// Seek sets the offset for the next Read or Write to offset,
// interpreted according to whence.
func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
