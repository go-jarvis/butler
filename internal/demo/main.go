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
	app = jarvis.NewApp(
		jarvis.WithName("Demo2s"),
	)
)

func init() {

	config := &struct {
		Server *Server
	}{
		Server: server,
	}

	app.Conf(config)

	server.engine.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

}

func main() {

	app.AddCommand("hello", func(args ...string) {
		hello()
	})
	// server.Run()
	// app.Run(server)

	// app.HelpCmd()
	app.Run(server)

}

func hello() {
	fmt.Println("hello go-jarvis")
}
