package main

import (
	"crypto/sha256"
	"github.com/awnumar/memguard"
	"os"
)

func calculateChecksum(filePath string) (*memguard.LockedBuffer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(data)
	return memguard.NewBufferFromBytes(hash[:]), nil
}
