package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"time"

	"bytes"
	"github.com/awnumar/memguard"
)

func init() {
	// Safely terminate if interrupted.
	slog.Info("Initiating interrupt-termination protection")
	memguard.CatchInterrupt()

	// Purge the session on return
	slog.Info("Initiating memguard purge on termination")
	defer memguard.Purge()
}

func secureCompareChecksum(currentChecksum *memguard.LockedBuffer) bool {
	initialChecksumBuf, err := initialChecksumEnclave.Open()
	if err != nil {
		fmt.Println("Error opening initial checksum:", err)
		return false
	}
	defer initialChecksumBuf.Destroy()

	return bytes.Equal(currentChecksum.Bytes(), initialChecksumBuf.Bytes())
}

func securePeriodicCheck(filePath string) {
	for {
		currentChecksum, err := calculateProtectedChecksum(filePath)
		if err != nil {
			slog.Error("Error calculating checksum:", "error", err)
			continue
		}
		currentChecksum.Destroy()

		if !secureCompareChecksum(currentChecksum) {
			slog.Info("Altered file detected")
		} else {
			slog.Info("File integrity intact")
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	testData, err := newTestData("/tmp/demo_dir", "demo_file.txt")
	if err != nil {
		slog.Error("Failed to create test data:", "error", err)
		os.Exit(1)
	}
	slog.Info("Test data has been created", "path", testData.path)

	initialChecksum, err := calculateProtectedChecksum(testData.path)
	if err != nil {
		slog.Error("Error calculating initial checksum:", "error", err)
		os.Exit(1)
	}
	defer initialChecksum.Destroy()

	slog.Info("Initial checksum ", "checksum", initialChecksum)

	slog.Info("Calling storeInitialChecksum")
	storeInitialChecksum(initialChecksum)
	os.Exit(0)

	slog.Info("Initial checksum of %q: %s", testData.path, hex.EncodeToString(initialChecksum.Bytes()))
	// Run the periodic check
	securePeriodicCheck(testData.path)
}
