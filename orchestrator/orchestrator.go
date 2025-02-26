package orchestrator

import (
	"finalTaskLMS/orchestrator/internal/api/handlers"
	"finalTaskLMS/orchestrator/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

var once sync.Once

const prefix string = "/api/v1"

type Server struct {
	R *gin.Engine
	O *types.Orchestrator
}

func NewOrchestratorServer() *Server {
	var s Server
	once.Do(func() {
		s = Server{
			O: types.NewOrchestrator(),
		}
	})
	return &s
}

// ConfigureRouter / специальная функция, которую можно менять для более гибкой конфигурации роутера библиотеки Gin
func (s *Server) ConfigureRouter() {
	r := gin.Default()

	r.POST(prefix+"/calculate", handlers.AddExpressionHandler(s.O))
	r.GET(prefix+"/expressions", handlers.GetExpressions(s.O))
	r.GET(prefix+"/expressions/:id", handlers.GetExpressionsById(s.O))

	s.R = r
}

func (s *Server) RunServer(int) {
	if err := s.R.Run(":8080"); err != nil {
		_ = fmt.Errorf("q")
	}
}
