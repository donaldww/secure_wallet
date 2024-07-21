package main

import (
	"crypto/sha256"
	"github.com/awnumar/memguard"
	"log/slog"
	"os"
)

var initialChecksumEnclave *memguard.Enclave

func calculateProtectedChecksum(filePath string) (*memguard.LockedBuffer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(data)
	return memguard.NewBufferFromBytes(hash[:]), nil
}

func storeInitialChecksum(checksum *memguard.LockedBuffer) {
	slog.Info("Storing initial checksum", "checksum", checksum.Bytes())
	os.Exit(0)
	initialChecksumEnclave = memguard.NewEnclave(checksum.Bytes())
}
