# butler

1. 根据配置 `config{}` 生成对应的 `default.yml` 配置文件。 
2. 读取依次配置文件 `default.yml, config.yml` + `分支配置文件.yml` + `环境变量`
    + 根据 GitlabCI, 分支配置文件 `config.xxxx.yml`
    + 如没有 CI, 读取本地文件: `local.yml`

## example

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

func (s *Server) SetDefaults() {
	if s.Port == 0 {
		s.Port = 80
	}
}

func (s *Server) Init() {
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