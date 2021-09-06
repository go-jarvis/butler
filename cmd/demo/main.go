package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gobutler/butler"
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

	app := butler.App{
		Name: "Server",
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
