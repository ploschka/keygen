package keygen

import (
	"crypto/rand"
	"io"
)

func GenerateRand(length uint64) ([]byte, error) {
	key := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
