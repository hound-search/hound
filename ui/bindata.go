package ui

import (
	"embed"
	"io"
	"os"
)

// content holds our static web server content.
//
//go:embed ui/.build/ui/**
var assetsFS embed.FS

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	f, err := assetsFS.Open("ui/.build/ui/" + name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	f, err := assetsFS.Open("ui/.build/ui/" + name)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}
