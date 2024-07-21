package main

import (
	"bytes"
	"github.com/awnumar/memguard"
	"log/slog"
	"time"
)

type Enclave struct {
	Enclave *memguard.Enclave
}

func NewEnclave() *Enclave {
	return &Enclave{}
}

func (e *Enclave) storeBuffer(checksum *memguard.LockedBuffer) {
	e.Enclave = checksum.Seal()
}

func (e *Enclave) runPeriodicCheck(filePath string, scanInterval time.Duration) {
	for {
		currentChecksum, err := calculateChecksum(filePath)
		if err != nil {
			slog.Error("Error calculating checksum", "error", err)
			continue
		}
		var fileStatus = "unknown"
		if !e.compareChecksum(currentChecksum) {
			fileStatus = "altered"
		} else {
			fileStatus = "intact"
		}
		slog.Info("File check", "status", fileStatus)
		time.Sleep(scanInterval * time.Second)
	}
}

func (e *Enclave) compareChecksum(currentChecksum *memguard.LockedBuffer) bool {
	initialChecksumBuf, err := e.Enclave.Open()
	if err != nil {
		slog.Error("Error opening enclave checksum", "error", err)
		return false
	}
	return bytes.Equal(currentChecksum.Bytes(), initialChecksumBuf.Bytes())
}
