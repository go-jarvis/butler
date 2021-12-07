package launcher

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Job interface {
	Name() string
	Run() error
}

type Shutdown interface {
	Shutdown() error
}

type Launcher struct {
	// 任务队列
	jobqueue chan Job
}

func (la *Launcher) Launch(jobs ...Job) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// lazy initial
	if la.jobqueue == nil {
		la.jobqueue = make(chan Job, len(jobs))
	}

	la.enqueue(jobs...)

	for {
		job, open := <-la.jobqueue
		if !open {
			logrus.Warn("jobqueue closed")
			break
		}

		go la.launch(ctx, job)
	}

}

func (la *Launcher) enqueue(jobs ...Job) {
	for _, job := range jobs {
		logrus.Debugf("%s 加入启动队列", job.Name())
		la.jobqueue <- job
	}
}

func (la *Launcher) launch(ctx context.Context, job Job) {

	logrus.Infof("启动程序: %s", job.Name())
	err := job.Run()

	if err != nil {
		logrus.Errorf("%s 程序运行失败: %v", job.Name(), err)
		// 重启
		la.enqueue(job)
	}
}
