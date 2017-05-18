package ops

import (
	"os"
	"io"
	"path/filepath"
)

const modeReadWrite os.FileMode = 0666

// LocalFile implements Storage on an OS-based file system
type LocalFile struct {
	root string
}

// NewLocalFile creates a new LocalFile object.
func NewLocalFile(root string) *LocalFile {
	fs := &LocalFile{root: root}
	return fs
}

// WriteTo reads key from the local filesystem and writes the bytes to w
func (fs *LocalFile) WriteTo(key string, w io.Writer) error {
	filename := fs.join(key)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

// ReadFrom reads from io.Reader r and writes the data to the local file system
func (fs *LocalFile) ReadFrom(key string, r io.Reader) error {
	filename := fs.join(key)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, modeReadWrite)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

// Delete remove the
func (fs *LocalFile) Delete(key string) error {
	filename := fs.join(key)
	return os.Remove(filename)
}

func (fs *LocalFile) join(elem ...string) string {
	args := append([]string{fs.root}, elem...)
	return filepath.Join(args...)
}

