package main

import (
	"log/slog"
	"os"

	"github.com/awnumar/memguard"
)

func init() {
	// Safely terminate if interrupted.
	slog.Info("Initiating interrupt protection")
	memguard.CatchInterrupt()

	// Purge the session on return
	slog.Info("Initiating memguard purge")
	defer memguard.Purge()
}

func main() {
	// Create a test file.
	testData, err := NewTestData("/tmp/demo_dir", "demo_file.txt")
	if err != nil {
		slog.Error("Failed to create test data", "error", err)
		os.Exit(1)
	}
	slog.Info("Test data created", "path", testData.Path)

	// Get the checksum of the test file as a memguard Lockfile.
	checksum, err := calculateChecksum(testData.Path)
	if err != nil {
		slog.Error("Error calculating initial checksum", "error", err)
		os.Exit(1)
	}

	// Store the initial checksum in a memguard enclave.
	enclave := NewEnclave()
	enclave.storeChecksum(checksum)

	// Rerun the checksum periodically, checking the new checksum
	// value against the stored checksum.
	enclave.securePeriodicCheck(testData.Path)
}
