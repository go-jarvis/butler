package jarvis

import (
	"fmt"
	"os"
	"testing"
	"text/template"
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
