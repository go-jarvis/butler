![](https://avatars.githubusercontent.com/u/82073077?s=400&u=f51fd1a2c01103122f249b4539fafb2495a109b1&v=4)

# jarivs

```bash
go get -u github.com/go-jarvis/jarvis
```

1. 根据配置 `config{}` 生成对应的 `default.yml` 配置文件。 
2. 读取依次配置文件 `default.yml, config.yml` + `分支配置文件.yml` + `环境变量`
    + 根据 GitlabCI, 分支配置文件 `config.xxxx.yml`
    + 如没有 CI, 读取本地文件: `local.yml`

## requeire

1. config 对象中的结构体中， 使用 `env:""` tag 才能的字段才会被解析到 **default.yml** 中。 也只有这些字段才能通过 **配置文件** 或 **环境变量** 进行初始化赋值。

2. config 中的对象需要有  `SetDefaults()` 和 `Init()` 方法。
    + `SetDefaults` 方法用于结构体设置默认值
    + `Init` 方法用于根据默认值初始化


## example

初始化代码如下

```go
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-jarvis/jarvis"
)

type Server struct {
	Listen string `env:"addr"`
	Port   int    `env:"port"`

	engine *gin.Engine
}

// SetDefaults values **Important**
func (s *Server) SetDefaults() {
	if s.Port == 0 {
		s.Port = 80
	}
}

// Initialize target, **Important**
func (s *Server) Initialize() {
	s.SetDefaults()

	if s.engine == nil {
		s.engine = gin.Default()
	}
}

func (s Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.Listen, s.Port)

	return s.engine.Run(addr)
}

func main() {
	server := &Server{}

	app := jarvis.App{
		Name: "Demo",
	}

	config := &struct {
		Server *Server
	}{
		Server: server,
	}
	// app.Save(config)

	app.Conf(config)
	// fmt.Println(config.Server.Port)

	server.Run()

}

```

### config

生成配置文件如下

```yaml
Demo__Server_addr: ""
Demo__Server_port: 80
```

在启动过程中， 如果环境变量中有同名变量, (例如 `Demo__Server_port`), 该变量值将被读取， 并复制给对应的字段。


## Launcher 启动器

jarvis Launcher 支持管理多任务并行启动，

如果服务满足接口 `IJob` 接口， 则可以通过 Launcher 管理启动和异常重启。

```go
type IJob interface {
	Name() string
	Run() error
}
```

如果服务同时满足 IShutdown 接口， 则可以 Launcher 可以通过 **信号** 或 **context** 触发程序安全退出。

```go
type IShutdown interface {
	Shutdown() error
}
```


### demo

demo [main.go](/internal/launcher/main.go)

```go

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 16*time.Second)
	defer cancel()

	la := &launcher.Launcher{}
	la.Launch(ctx, &Runner2{})
}
```

### jarvis appctx 支持 cobra.Command 命令模式

支持添加子命令

```go
func main() {

	app.AddCommand("hello", func(args ...string) {
		// target = strings.Join(args, ", ")
		hello()
		helloFailed(target)
	},
		func(cmd *cobra.Command) {
			cmd.Long = "say hello"
		},
	)

	app.Run(server)

}


// say hello
// Usage:
//   Demo2s hello [flags]
// Flags:
//   -h, --help            help for hello
```