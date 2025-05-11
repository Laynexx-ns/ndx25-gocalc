package agent

import (
	"finalTaskLMS/internal/services/agent/internal/handlers"
	"finalTaskLMS/internal/services/agent/types"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	"strconv"
	"sync"
)

var once sync.Once

type Server struct {
	R *echo.Echo
	A *types.Agent
}

func NewAgentServer() *Server {
	var s Server
	once.Do(func() {
		s = Server{
			A: &types.Agent{},
		}
	})
	return &s
}

func (s *Server) InitializeAgent() {
	res, err := strconv.Atoi(os.Getenv("LIMIT_OF_GOROUTINES"))
	if err != nil {
		fmt.Println("err : can't parse LIMIT_OF_GOROUTINES (env)")
	}
	s.A.LimitOfGoroutines = res
	s.A.CountOfGoroutines = 0
}

func (s *Server) ConfigureRouter() {
	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "( ´ ꒳ ` )")
	})

	e.Use(middleware.CORS())

	s.R = e
}

func (s *Server) RunServer() {
	go handlers.CycleTask(s.A)

	if err := s.R.Start(":8081"); err != nil {
		log.Fatal("can[k")
	}

}
