package jarvis

import (
	_ "embed"
	"os"
	"text/template"

	"github.com/sirupsen/logrus"
)

//go:embed tmpl/Dockerfile.tmpl
var dockerfileTmpl string

func (app *AppCtx) dockerizeCommand() {
	tmpl, _ := template.New("dockerfile").Parse(dockerfileTmpl)

	fobj, err := os.OpenFile("Dockerfile.default", os.O_TRUNC|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		logrus.Errorf("create Dockerfile.default failed: %v", err)
	}
	defer fobj.Close()

	data := struct {
		Name     string
		WorkRoot string
	}{
		Name:     app.name,
		WorkRoot: "internal/demo",
	}

	err = tmpl.Execute(fobj, data)
	if err != nil {
		logrus.Errorf("write Dockerfile.default failed: %v", err)
	}
}
