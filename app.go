package jarvis

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-jarvis/jarvis/launcher"
	"github.com/tangx/envutils"
)

// App 配置文件管理器
type App struct {
	name string
}

// NewApp 创建一个配置文件管理器
func NewApp() *App {
	return &App{}
}

// setDefaults 设置默认值
func (app *App) setDefaults() {
	if app.name == "" {
		app.name = "APP"
	}
}

// WithName 为 App 指定一个名字
func (app *App) WithName(name string) *App {
	app.name = name
	return app
}

// Conf 解析配置， 并在 config 目录下生成 xxx.yml 文件
func (app *App) Conf(config interface{}) error {
	app.setDefaults()

	// call SetDefaults
	if err := envutils.CallSetDefaults(config); err != nil {
		return err
	}

	// write config
	data, err := envutils.Marshal(config, app.name)
	if err != nil {
		return err
	}
	_ = os.MkdirAll("./config", 0755)
	_ = os.WriteFile("./config/default.yml", data, 0644)

	// load config from files
	for _, _conf := range []string{"default.yml", "config.yml", refConfig()} {
		_conf := filepath.Join("./config/", _conf)
		err = envutils.UnmarshalFile(config, app.name, _conf)
		if err != nil {
			log.Println(err)
		}
	}

	// load config from env
	err = envutils.UnmarshalEnv(config, app.name)
	if err != nil {
		log.Print(err)
	}

	// CallInit
	if err := envutils.CallInit(config); err != nil {
		return err
	}

	return nil
}

// refConfig 根据 gitlab ci 环境变量创建与分支对应的配置文件
func refConfig() string {
	// gitlab ci
	ref := os.Getenv("CI_COMMIT_REF_NAME")

	if len(ref) != 0 {
		return refFilename(ref)
	}

	return `local.yml`
}

// refFilename 根据 ref 信息返回对应的配置文件名称
func refFilename(ref string) string {
	// feat/xxxx
	parts := strings.Split(ref, "/")
	feat := parts[len(parts)-1]               // xxxx
	return fmt.Sprintf("config.%s.yml", feat) // config.xxxx.yml
}

// Run 启动服务
func (app *App) Run(jobs ...launcher.IJob) {
	ctx := context.Background()
	app.RunContext(ctx, jobs...)
}

// RunContext 启动服务
func (app *App) RunContext(ctx context.Context, jobs ...launcher.IJob) {

	launcher := &launcher.Launcher{}
	launcher.Launch(ctx, jobs...)
}
