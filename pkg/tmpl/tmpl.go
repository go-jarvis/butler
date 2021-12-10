package tmpl

import (
	"embed"
	"io/fs"
	"path/filepath"
)

const (
	PlaceHolder_ProjectName = `jarvis-demo`
	PlaceHolder_FileSuffix  = ".tmpl"
)
const (
	_PATH_Templates        = "templates"
	_PATH_TemplatesProject = "templates/project"
)

//go:embed templates
var dir embed.FS

func GetFile(filename string) ([]byte, error) {
	target := filepath.Join(_PATH_Templates, filename)
	return dir.ReadFile(target)
}

func ReadDir(path string) ([]fs.DirEntry, error) {
	return dir.ReadDir(path)
}

func ReadProjectDir(path string) ([]fs.DirEntry, error) {
	fullpath := filepath.Join(_PATH_TemplatesProject, path)
	return ReadDir(fullpath)
}

func ReadProjectFile(filename string) ([]byte, error) {

	target := filepath.Join(_PATH_TemplatesProject, filename)
	return dir.ReadFile(target)
}
