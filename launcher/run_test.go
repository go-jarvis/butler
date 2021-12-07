package launcher

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type Runner2 struct {
}

func (r *Runner2) Name() string {
	return "runner2s"
}

func (r *Runner2) Shutdown() error {
	return nil
}

func (r *Runner2) Run() error {
	time.Sleep(3 * time.Second)
	return fmt.Errorf("sleep 5 s")
}

func Test_Run(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)

	la := &Launcher{}
	la.Launch(&Runner2{})
}
