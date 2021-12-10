package tmpl

import (
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed templates/**
var dir embed.FS

func GetFile(filename string) ([]byte, error) {

	entries, _ := dir.ReadDir("templates")

	for _, entry := range entries {
		if entry.Name() == filename {
			f := filepath.Join("templates", filename)
			return dir.ReadFile(f)
		}
	}

	return nil, fmt.Errorf("%s not found", filename)
}
