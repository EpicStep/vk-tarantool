package vk_tarantool

import (
	"embed"
	"io/fs"
)

//go:embed frontend
var embedFrontend embed.FS

func GetFrontendAssets() (fs.FS, error) {
	f, err := fs.Sub(embedFrontend, "frontend")
	if err != nil {
		return nil, err
	}

	return f, nil
}