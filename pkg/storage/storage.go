/*
Package storage manages content stores used to save/get/delete files using their blake2b hash value.
*/
package storage

import (
	"fmt"
	"io"

	"golang.org/x/crypto/blake2b"
)

// HashSize is the byte size of the blake2b hash.
const HashSize = 32

// Store represents a storage in which media can be managed through its blake2b hash.
type Store interface {
	SaveFile(file io.ReadSeeker) (string, error)
	DeleteFileWithHash(hash string) error
	GetFileWithHash(hash string) (io.ReadSeeker, error)
}

func getHash(reader io.ReadSeeker) (string, error) {
	hash, err := blake2b.New(HashSize, nil)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(hash, reader)
	if err != nil {
		return "", err
	}

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
