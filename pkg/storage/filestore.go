package storage

import (
	"io"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
)

const HashSize = 32

var fs = afero.NewOsFs()

// FileStore represents a store using the local filesystem.
type FileStore struct {
	Path string
}

// NewFileStore creates store at given path.
func NewFileStore(path string) (*FileStore, error) {
	cleanPath := filepath.Clean(path)
	err := fs.MkdirAll(cleanPath, 0755)
	if err != nil {
		return nil, err
	}

	return &FileStore{
		Path: cleanPath,
	}, nil
}

// Save writes the file into the store. The file's hash will be used
// as filename. It returns the file's hash and any error encountered.
func (s FileStore) Save(file io.ReadSeeker) (string, error) {
	hash, err := getHash(file)
	if err != nil {
		return "", err
	}

	exists, err := afero.Exists(fs, path.Join(s.Path, hash))
	if err != nil {
		return "", err
	}
	if exists {
		return hash, nil
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", nil
	}

	if err = s.write(hash, file); err != nil {
		return "", err
	}

	return hash, nil
}

func (s FileStore) write(hash string, file io.Reader) error {
	newFile, err := fs.Create(path.Join(s.Path, hash))
	if err != nil {
		return err
	}

	_, err = io.Copy(newFile, file)
	if err != nil {
		return err
	}

	if err := newFile.Close(); err != nil {
		return err
	}

	return nil
}

// Get returns a io.ReadSeeker from the file with given hash in the store.
func (s FileStore) Get(hash string) (io.ReadSeeker, error) {
	exists, err := afero.Exists(fs, path.Join(s.Path, hash))
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, nil
	}

	file, err := fs.Open(path.Join(s.Path, hash))
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Delete deletes the file with given hash from the store.
func (s FileStore) Delete(hash string) error {
	return fs.Remove(path.Join(s.Path, hash))
}
