package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-jarvis/jarvis/launcher"
	"github.com/sirupsen/logrus"
)

type Runner2 struct {
}

func (r *Runner2) Appname() string {
	return "runner2s"
}

func (r *Runner2) Shutdown(ctx context.Context) error {
	return nil
}

func (r *Runner2) Run() error {
	time.Sleep(3 * time.Second)
	return fmt.Errorf("sleep 5 s")
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 16*time.Second)
	defer cancel()

	la := &launcher.Launcher{}
	la.Launch(ctx, &Runner2{})
}
