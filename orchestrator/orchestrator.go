package orchestrator

import (
	"finalTaskLMS/orchestrator/internal/api/handlers"
	"finalTaskLMS/orchestrator/internal/models"
	"finalTaskLMS/orchestrator/pkg/calc"
	"finalTaskLMS/orchestrator/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
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

func (s *Server) ConfigureRouter() {
	r := gin.Default()

	r.POST(prefix+"/calculate", AddExpressionHandler(s))
	r.GET(prefix+"/expressions", handlers.GetExpressions(s.O))
	r.GET(prefix+"/expressions/:id", handlers.GetExpressionsById(s.O))
	r.GET(prefix+"/queue", func(c *gin.Context) {
		c.JSON(200, s.O.Queue)
	})
	r.GET("internal/task", handlers.GetTasks(s.O))
	r.POST("internal/task", handlers.GetExpressionResult(s.O))

	s.R = r
}

func (s *Server) RunServer(int) {

	if err := s.R.Run(":8080"); err != nil {
		_ = fmt.Errorf("q")
	}
}

func (s *Server) CreateTasks() {

	go func() {
		for i, v := range s.O.Expressions {
			go func(i int, v models.Expressions) {
				s.O.Mu.Lock()
				s.O.Chans[i] = make(chan float64, 1)
				s.O.Errchans[i] = make(chan error, 1)
				s.O.Mu.Unlock()

				if v.Status == "pending" {
					calc.Calc(v.Expression, s.O.Chans[i], s.O.Errchans[i], i, s.O)

					select {
					case res := <-s.O.Chans[i]:
						s.O.Mu.Lock()
						for j, e := range s.O.Expressions {
							if e.Id == i {
								s.O.Expressions[j].Status = "successfully calculated"
								s.O.Expressions[j].Result = res
							}
						}
						s.O.Mu.Unlock()
					case <-s.O.Errchans[i]:
						close(s.O.Chans[i])
						s.O.Mu.Lock()
						for j, e := range s.O.Expressions {
							if e.Id == i {
								s.O.Expressions[j].Status = "failed"
								s.O.Expressions[j].Result = 0
							}
						}
						s.O.Mu.Unlock()
					}
				}
			}(i, v)
		}
	}()
}

func AddExpressionHandler(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.JSON(405, gin.H{
				"error": "this method not allowed",
			})
		}

		var expression models.UserExpressions
		if err := c.ShouldBindJSON(&expression); err != nil {
			c.JSON(400, gin.H{
				"error": "invalid JSON",
			})
			return
		}

		valid, err := regexp.MatchString("^[0-9)(*/+-]+$", expression.Expression)
		if err != nil || !valid {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "request contains forbidden symbols",
			})
			return
		}

		s.O.Mu.Lock()
		defer s.O.Mu.Unlock()
		s.O.Expressions = append(s.O.Expressions, models.Expressions{
			Id:         len(s.O.Expressions),
			Status:     "pending",
			Result:     0,
			Expression: expression.Expression,
		})
		c.JSON(200, gin.H{
			"id": s.O.Expressions[len(s.O.Expressions)-1].Id,
		})
		s.CreateTasks()

	}
}
