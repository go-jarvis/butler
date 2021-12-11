package appctx

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-jarvis/jarvis/pkg/launcher"
	"github.com/spf13/cobra"
	"github.com/tangx/envutils"
)

// New 创建一个配置文件管理器
func New() *AppCtx {
	app := &AppCtx{}
	app.cmd = &cobra.Command{}
	return app
}

func (app *AppCtx) WithOptions(opts ...AppCtxOption) *AppCtx {
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// AppCtx 配置文件管理器
type AppCtx struct {
	name     string
	rootdir  string
	helpMode bool
	cmd      *cobra.Command
}

type AppCtxOption = func(app *AppCtx)

func WithHelpMode() AppCtxOption {
	return func(app *AppCtx) {
		app.helpMode = true
	}
}

// WithRoot 指定项目根目录， 即 go.mod 所在的目录。
//  root 可以是相对 main.go 的相对路径， 也是可以是绝对路径
func WithRoot(root string) AppCtxOption {
	return func(app *AppCtx) {
		abs, _ := filepath.Abs(root)
		app.rootdir = abs
	}
}

// WithName 设置 name
func WithName(name string) AppCtxOption {
	return func(app *AppCtx) {
		app.name = name
	}
}

// setDefaults 设置默认值
func (app *AppCtx) setDefaults() {
	if app.name == "" {
		app.name = "app"
	}

	// if rootdir wat not set, use current path as rootdir.
	if app.rootdir == "" {
		app.rootdir, _ = os.Getwd()
	}
}

// Conf 解析配置， 并在 config 目录下生成 xxx.yml 文件
func (app *AppCtx) Conf(config interface{}) error {
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

// AddCommand 添加子命令
// ex: AddCommand(migrate, module.Migrate())
//
// cmdOpts can be flags options
// ex:
//
// func WithFlags(flag string) func(*cobra.Command) {
// 	return func(cmd *cobra.Command) {
// 		cmd.Flags().StringVarP(&flag, "targets", "t", "nothing", "specify targets")
// 	}
// }
func (app *AppCtx) AddCommand(use string, fn func(args ...string), cmdOpts ...func(*cobra.Command)) {
	subCmd := &cobra.Command{
		Use: use,
	}

	subCmd.Run = func(cmd *cobra.Command, args []string) {

		if app.helpMode {
			_ = cmd.Help()
		}

		fn(args...)
	}

	for _, opt := range cmdOpts {
		opt(subCmd)
	}

	app.cmd.AddCommand(subCmd)
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
func (app *AppCtx) Run(jobs ...launcher.IJob) {
	ctx := context.Background()
	app.RunContext(ctx, jobs...)
}

// RunContext 启动服务
func (app *AppCtx) RunContext(ctx context.Context, jobs ...launcher.IJob) {

	launcher := &launcher.Launcher{}

	app.cmd.Use = app.name

	// 添加命令
	app.cmd.Run = func(cmd *cobra.Command, args []string) {

		if app.helpMode {
			_ = cmd.Help()
		}

		launcher.Launch(ctx, jobs...)
	}

	// dockerize
	app.AddCommand("dockerize", func(args ...string) {
		app.dockerizeCommand()
	}, func(c *cobra.Command) {
		c.Short = "dockerize project"
	})

	// 启动服务
	if err := app.cmd.Execute(); err != nil {
		panic(err)
	}
}
