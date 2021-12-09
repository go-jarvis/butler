package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-jarvis/jarvis"
	"github.com/spf13/cobra"
)

type Server struct {
	Listen string `env:"addr"`
	Port   int    `env:"port"`

	engine *gin.Engine
}

func (s *Server) SetDefaults() {
	if s.Port == 0 {
		s.Port = 8098
	}
}

func (s *Server) Init() {
	s.SetDefaults()

	if s.engine == nil {
		s.engine = gin.Default()
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.Listen, s.Port)

	return s.engine.Run(addr)
}

func (s *Server) Appname() string {
	return "http-serserver"
}

var (
	server = &Server{}

	// app := jarvis.NewApp().WithName("Demo2")
	app = jarvis.New(
		jarvis.WithName("Demo2s"),
		jarvis.WithHelpMode(),
	)
)

func init() {

	config := &struct {
		Server *Server
	}{
		Server: server,
	}

	_ = app.Conf(config)

	server.engine.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

}

var target = ""

func main() {

	cmdopt := func(cmd *cobra.Command) {
		// fmt.Println("添加 targets flag")
		cmd.Flags().StringVarP(&target, "target", "t", "defualt value", "say hello to targets")
	}

	/*
		todo: helloFailed: 无法获取 target 的值
		solution: 因为 cmd options 是作为一个参数在 app.AddCommand 中执行的。 类似 defer， 在 压栈??? 的时候， target 的值就已经确定了， 并在 subcommand 中被应用。。 因此后面在 cobra 解析 flag 的时候， 无法获取到新值。 对于这种值， 可以使用全局变量（比较丑陋） 或者 引用类型（指针）作为参数传递。
	*/

	// var target = "" // 错误地点
	app.AddCommand("hello", func(args ...string) {
		// target = strings.Join(args, ", ")
		hello()
		helloFailed(target)
	},
		cmdopt,

		func(cmd *cobra.Command) {
			// fmt.Println("添加 targets flag")
			cmd.Long = "say hello"
		},
	)

	app.Run(server)

}

func hello() {
	if len(target) == 0 {
		fmt.Println("hello go-jarvis")
		return
	}

	fmt.Println("hello", target)
}

func helloFailed(target string) {
	if len(target) == 0 {
		fmt.Println("helloFailed go-jarvis")
		return
	}

	fmt.Println("helloFailed", target)
}
