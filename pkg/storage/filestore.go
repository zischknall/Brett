package storage

import (
	"io"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

// FileStore implements a store using the local filesystem.
type FileStore struct {
	Path string
}

// GetFileStore creates a FileStore at given path.
func GetFileStore(path string) (*FileStore, error) {
	cleanPath := filepath.Clean(path)
	err := fs.MkdirAll(cleanPath, 0755)
	if err != nil {
		return nil, err
	}

	return &FileStore{
		Path: cleanPath,
	}, nil
}

// SaveFile saves the file into the FileStore.
// Returns the hash of file.
func (s FileStore) SaveFile(file io.ReadSeeker) (string, error) {
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

	if err = s.createFileWithHash(hash, file); err != nil {
		return "", err
	}

	return hash, nil
}

func (s FileStore) createFileWithHash(hash string, file io.Reader) error {
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

// GetFileWithHash returns an io.ReadSeeker from the file with given hash in the FileStore.
func (s FileStore) GetFileWithHash(hash string) (io.ReadSeeker, error) {
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

// DeleteFileWithHash deletes the file with given hash from the FileStore.
func (s FileStore) DeleteFileWithHash(hash string) error {
	return fs.Remove(path.Join(s.Path, hash))
}
