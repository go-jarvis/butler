package launcher

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	enqueueEnabled = true
)

type IJob interface {
	Name() string
	Run() error
}

type IShutdown interface {
	Shutdown(ctx context.Context) error
}

type Launcher struct {
	// 任务队列
	jobqueue chan IJob
	// 重启次数
	jobs map[IJob]int
}

// Launch 启动程序
func (la *Launcher) Launch(ctx context.Context, jobs ...IJob) {

	la.lazyInitial(len(jobs))
	la.enqueue(ctx, jobs...)

	defer func() {
		la.shutdown(ctx)
	}()

	// 监听捕获信号
	// https://colobu.com/2015/10/09/Linux-Signals/
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigs:
			// 发送关闭信号
			logrus.Infof("catch exit signal: %v", sig)
			la.close()

			return
		case <-ctx.Done():
			logrus.Infof("context done: %v", ctx.Err())

			la.close()
			return

		case job, open := <-la.jobqueue:
			// 执行任务
			// 如果 jobqueue 不关闭， 将一直阻塞
			time.Sleep(1 * time.Second)
			if !open {
				logrus.Warn("jobqueue closed")
				return
			}
			go la.launch(ctx, job)

		}
	}
}

// lazyInitial 懒加载初始化
func (la *Launcher) lazyInitial(n int) {

	if la.jobqueue == nil {
		la.jobqueue = make(chan IJob, n)

	}

	if la.jobs == nil {
		la.jobs = make(map[IJob]int, n)
	}
}

// close 关闭资源
func (la *Launcher) close() {

	close(la.jobqueue)
	enqueueEnabled = false
}

// enqueue 任务入队
func (la *Launcher) enqueue(ctx context.Context, jobs ...IJob) {

	for _, job := range jobs {
		logrus.Debugf("job %s enqueue", job.Name())

		// counter++
		la.jobs[job] += 1

		la.jobqueue <- job
	}
}

// launch 启动任务
func (la *Launcher) launch(ctx context.Context, job IJob) {

	// 捕获程序内部 panic
	defer func() {
		if err := recover(); err != nil {
			// 重启
			logrus.Errorf("job %s runs failed: %v", job.Name(), err)
			if enqueueEnabled {
				la.enqueue(ctx, job)
			}
		}
	}()

	logrus.Infof("job %s (re)start at %d times", job.Name(), la.jobs[job])

	err := job.Run()
	panic(err)

}

// shutdown 优雅关闭任务
func (la *Launcher) shutdown(ctx context.Context) {
	logrus.Info("START to STUT all jobs DOWN: ")

	for job := range la.jobs {

		app, ok := job.(IShutdown)
		if !ok {
			logrus.Warnf("%s has NO SHUTDOWN method: skip", job.Name())
			continue
		}

		logrus.Infof("shutting down: %s", job.Name())
		err := app.Shutdown(ctx)
		if err != nil {
			logrus.Errorf("%s shutdown failed: %v", job.Name(), err)
		}
	}
}
