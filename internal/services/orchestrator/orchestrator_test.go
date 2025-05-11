package orchestrator

import (
	"bytes"
	"encoding/json"
	"finalTaskLMS/internal/models"
	"finalTaskLMS/internal/services/orchestrator/types"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestAddExpressionHandler_ValidExpression(t *testing.T) {
	s := NewOrchestratorServer()
	s.O = types.NewOrchestrator()
	s.ConfigureRouter()

	expr := models.UserExpressions{Expression: "2+2"}
	jsonValue, _ := json.Marshal(expr)
	request, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	s.R.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.Code)
	}
}

func TestAddExpressionHandler_InvalidExpression(t *testing.T) {
	s := NewOrchestratorServer()
	s.O = types.NewOrchestrator()
	s.ConfigureRouter()

	expr := models.UserExpressions{Expression: "2a+2"}
	jsonValue, _ := json.Marshal(expr)
	request, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	s.R.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", response.Code)
	}
}

func TestAddExpressionHandler_EmptyExpression(t *testing.T) {
	s := NewOrchestratorServer()
	s.O = types.NewOrchestrator()
	s.ConfigureRouter()

	expr := models.UserExpressions{Expression: ""}
	jsonValue, _ := json.Marshal(expr)
	request, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	s.R.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", response.Code)
	}
}

func TestAddExpressionHandler_Concurrency(t *testing.T) {
	s := NewOrchestratorServer()
	s.O = types.NewOrchestrator()
	s.ConfigureRouter()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			expr := models.UserExpressions{Expression: "2+2"}
			jsonValue, _ := json.Marshal(expr)
			request, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonValue))
			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			s.R.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", response.Code)
			}
		}()
	}
	wg.Wait()
}
