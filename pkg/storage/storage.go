/*
storage manages a content store indexed by file hashes.
*/
package storage

import (
	"fmt"
	"io"

	"golang.org/x/crypto/blake2b"
)

type Store interface {
	Save(file io.ReadSeeker) (string, error)
	Delete(hash string) error
	Get(hash string) (io.ReadSeeker, error)
}

func getHash(reader io.Reader) (string, error) {
	hash, err := blake2b.New(HashSize, nil)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(hash, reader)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
