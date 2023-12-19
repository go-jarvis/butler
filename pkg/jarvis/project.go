package jarvis

// 好像其实没什么用
// 当初是为了模仿 cobra 的方式， 创建一个项目初始化的命令。

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

// CreateProject create a new project
func (info *ProjectInfo) CreateProject() {

	info.initial()

	rootdir := "."
	entries, err := tmpl.ReadProjectDir(rootdir)
	if err != nil {
		panic(err)
	}

	info.walk(entries, rootdir)

}

// walk templates/project folder recursive
func (info *ProjectInfo) walk(entries []fs.DirEntry, dirname string) {
	err := info.render(dirname, true)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		name := entry.Name()

		fullpath := filepath.Join(dirname, name)

		// create file
		if !entry.IsDir() {
			err := info.render(fullpath, false)
			if err != nil {
				panic(err)
			}

			continue
		}

		// walkdir
		subEntries, err := tmpl.ReadProjectDir(fullpath)
		if err != nil {
			panic(err)
		}

		info.walk(subEntries, fullpath)
	}
}

// target calculate target
func (info *ProjectInfo) target(source string) string {
	// replace placeholder
	target := strings.ReplaceAll(source, tmpl.PlaceHolder_ProjectName, info.Name)
	target = strings.TrimSuffix(target, tmpl.PlaceHolder_FileSuffix)
	// join real path
	return filepath.Join(info.Workdir, target)

}

// render 创建目录或渲染文件
//
//	source 是 project 下文件的相对路径
func (info *ProjectInfo) render(source string, isDir bool) error {
	target := info.target(source)

	// create folder
	if isDir {
		err := os.MkdirAll(target, os.ModePerm)
		return err
	}

	// read file
	content, _ := tmpl.ReadProjectFile(source)

	fobj, err := os.OpenFile(target, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fobj.Close()

	t, _ := template.New(target).Parse(string(content))
	err = t.Execute(fobj, info)
	if err != nil {
		return err
	}

	return nil
}
