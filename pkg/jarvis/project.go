package jarvis

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-jarvis/jarvis/pkg/tmpl"
)

type ProjectInfo struct {
	Name    string `flag:"name" usage:"project name"`
	Workdir string `flag:"dir" usage:"project src path, \ndefault: <project_name>"`
	PkgName string `flag:"pkg" usage:"go module name, \ndefault: github.com/go-jarvis/<project_name>"`
}

// var makefile = tmpl.GetFile("Makefile")
var Project *ProjectInfo

func init() {
	Project = &ProjectInfo{
		Name: "app",
		// Workdir: "app",
	}
	// Project.PkgName = fmt.Sprintf("github.com/go-jarvis/%s", Project.Name)
}

func (info *ProjectInfo) initial() {
	if info.Workdir == "" {
		info.Workdir = info.Name
	}

	if info.PkgName == "" {
		info.PkgName = fmt.Sprintf("github.com/go-jarvis/%s", info.Name)
	}

}

// CreateProject 初始化项目
// {{ Workdir }}/cmd/{{ Name }}/main.go
// {{ Workdir }}/Makefile
// {{ Workdir }}/go.mod
func (info *ProjectInfo) CreateProject() {

	info.initial()

	for _, dir := range []string{
		info.Workdir,
		fmt.Sprintf("%s/cmd/%s", info.Workdir, info.Name),
	} {
		mkdir(dir)
	}

	for _, file := range []string{
		"Makefile",
	} {
		touch(file, info.Workdir)
	}
}

func mkdir(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
}

func touch(file string, dir string) {
	content, err := tmpl.GetFile(file)
	if err != nil {
		panic(err)
	}

	f := filepath.Join(dir, file)
	_ = ioutil.WriteFile(f, content, os.ModePerm)
}
