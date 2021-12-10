package appctx

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func Test_refConfig(t *testing.T) {

	for _, ref := range []string{
		"master",
		"develop",
		"feat/xxxx",
	} {
		fmt.Println(refFilename(ref))
	}
}

func TestTemplate(t *testing.T) {

	app := struct {
		name string
	}{
		name: "demo",
	}

	// readfile
	// tmpl, err := template.ParseFiles("tmpl/Dockerfile.tmpl")

	// read string
	tmpl, err := template.New("dockerfile").Parse(dockerfileTmpl)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name     string
		WorkRoot string
	}{
		Name:     app.name,
		WorkRoot: "internal/demo",
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

}

func Test_FilepathAbs(t *testing.T) {
	datas := []struct {
		root   string
		wanted string
	}{
		{root: "", wanted: "."},
		{root: ".", wanted: "."},
		{root: "..", wanted: "jarvis"},
		{root: "../../", wanted: "go-jarvis/jarvis"},
		{root: "../../..", wanted: "github.com/go-jarvis/jarvis"},
	}
	for _, data := range datas {

		t.Run(data.root, func(t *testing.T) {
			r := abs(data.root)
			sub := workdir(r)
			NewWithT(t).Expect(sub).To(Equal(data.wanted))
		})

	}
}

func abs(root string) string {
	p, err := filepath.Abs(root)
	if err != nil {
		logrus.Fatal(err)
	}

	return p
}
