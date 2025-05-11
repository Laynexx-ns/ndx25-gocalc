package orchestrator

import (
	"context"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal/api/handlers"
	"ndx/internal/services/orchestrator/pkg/calc"
	"ndx/internal/services/orchestrator/types"

	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"regexp"
	"sync"
	"sync/atomic"
	"time"
)

var once sync.Once

const prefix string = "/api/v1"

type Server struct {
	Server            *echo.Echo
	Orchestrator      *types.Orchestrator
	expressionCounter uint64
}

func NewOrchestratorServer() *Server {
	var s Server

	once.Do(func() {
		s = Server{
			Orchestrator: types.NewOrchestrator(),
		}
	})
	return &s
}
func (s *Server) ConfigureRouter() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	r.POST(prefix+"/calculate", AddExpressionHandler(s))
	r.GET(prefix+"/expressions", handlers.GetExpressions(s.O))
	r.GET(prefix+"/expressions/:id", handlers.GetExpressionsById(s.O))
	r.GET(prefix+"/queue", func(c *gin.Context) {
		c.JSON(200, s.Orchestrator.Queue)
	})
	r.GET("internal/task", handlers.GetTasks(s.Orchestrator))
	r.POST("internal/task", handlers.SetExpressionResult(s.Orchestrator))

	s.R = r
}

func (s *Server) RunServer(int) {

	if err := s.R.Run(":8080"); err != nil {
		_ = fmt.Errorf("q")
	}
}

func (s *Server) CreateTasks() {
	s.O.Mu.Lock()
	defer s.O.Mu.Unlock()

	for i, v := range s.O.Expressions {
		if v.Status != "pending" {
			continue
		}

		s.O.Expressions[i].Status = "processing"

		go func(id int, expr models.Expressions) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			resChan := make(chan float64, 1)
			errChan := make(chan error, 1)

			go func() {
				calc.Calc(expr.Expression, resChan, errChan, id, s.O)
			}()

			select {
			case res := <-resChan:
				s.O.Mu.Lock()
				for j, e := range s.O.Expressions {
					if e.Id == id {
						s.O.Expressions[j].Status = "successfully calculated"
						s.O.Expressions[j].Result = res
						break
					}
				}
				s.O.Mu.Unlock()

			case <-errChan:
				s.O.Mu.Lock()
				for j, e := range s.O.Expressions {
					if e.Id == id {
						s.O.Expressions[j].Status = "failed"
						s.O.Expressions[j].Result = 0
						break
					}
				}
				s.O.Mu.Unlock()

			case <-ctx.Done():
				s.O.Mu.Lock()
				for j, e := range s.O.Expressions {
					if e.Id == id {
						s.O.Expressions[j].Status = "timeout"
						break
					}
				}
				s.O.Mu.Unlock()
			}
		}(v.Id, v)
	}
}

func AddExpressionHandler(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var expression models.UserExpressions
		if err := c.ShouldBindJSON(&expression); err != nil {
			c.JSON(400, gin.H{"error": "invalid JSON"})
			return
		}

		valid, _ := regexp.MatchString("^[0-9)(*/+-]+$", expression.Expression)
		if !valid {
			c.JSON(422, gin.H{"error": "invalid characters"})
			return
		}

		s.O.Mu.Lock()
		newID := atomic.AddUint64(&s.expressionCounter, 1)
		newExpr := models.Expressions{
			Id:         int(newID),
			Status:     "pending",
			Result:     0,
			Expression: expression.Expression,
		}
		s.O.Expressions = append(s.O.Expressions, newExpr)
		s.O.Mu.Unlock()

		c.JSON(200, gin.H{"id": newID})
		s.CreateTasks()
	}
}
