package main

import (
	"os"
	"path/filepath"
)

type TestData struct {
	Path string
}

func NewTestData(directory string, filename string) (*TestData, error) {
	path, err := createInitialFile(directory, filename)
	if err != nil {
		return nil, err
	}
	return &TestData{
		Path: path,
	}, nil
}

func createInitialFile(directory, filename string) (string, error) {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(directory, filename)
	if err := os.WriteFile(path, []byte("Initial content"), 0644); err != nil {
		return "", err
	}
	return path, nil
}
