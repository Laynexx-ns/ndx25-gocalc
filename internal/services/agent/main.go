package agent

import (
	"github.com/labstack/echo/v4"
	"ndx/internal/services/agent/internal/config"
	"ndx/pkg/config"
	"sync"
)

var once sync.Once

type Server struct {
	Config config.Config
	Echo   *echo.Echo
	Agent  *config.Agent
}

func main() {

}
