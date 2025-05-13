package calc

import (
	"errors"
	"fmt"
	"log"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal/repo"
	"ndx/pkg/logger"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func isNum(r rune) bool {
	return unicode.IsDigit(r)
}

var priority = map[rune]int{
	'+': 1,
	'-': 1,
	'/': 2,
	'*': 2,
}

var resultChans = make(map[int]chan float64)

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/'
}

func EvaluateSimpleExpression(a, b float64, operand string, parentId int, repo *repo.TasksRepository) (float64, error) {

	id, err := repo.SavePrimeEvaluation(models.PrimeEvaluation{
		ParentID:      parentId,
		Arg1:          a,
		Arg2:          b,
		Operation:     operand,
		OperationTime: 0,
	})
	log.Println(id)
	res, err := WaitForEvaluationResult(repo, id, 5*time.Second)

	if err != nil {
		return 0, err
	}

	log.Println("res: ", res)
	return res, nil
}

func WaitForEvaluationResult(repo *repo.TasksRepository, Id int, timeout time.Duration) (float64, error) {
	start := time.Now()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s, err := repo.GetPrimeEvaluationByID(Id)
			logger.L().Logf(0, "prime eval: %v", s)

			if err != nil {
				return 0, err
			}
			if s.CompletedAt != nil {
				return s.Result, nil
			}
			if time.Since(start) > timeout {
				return 0, errors.New("evaluation timeout")
			}
		}
	}
}

func Parse(expression string) ([]string, error) {
	var result []string
	var operators []rune
	var num strings.Builder

	for _, ch := range expression {
		if isNum(ch) || ch == '.' {
			num.WriteRune(ch)
		} else {
			if num.Len() > 0 {
				result = append(result, num.String())
				num.Reset()
			}
			if isOperator(ch) {
				for len(operators) > 0 && priority[operators[len(operators)-1]] >= priority[ch] {
					result = append(result, string(operators[len(operators)-1]))
					operators = operators[:len(operators)-1]
				}
				operators = append(operators, ch)
			} else if ch == '(' {
				operators = append(operators, ch)
			} else if ch == ')' {
				for len(operators) > 0 && operators[len(operators)-1] != '(' {
					result = append(result, string(operators[len(operators)-1]))
					operators = operators[:len(operators)-1]
				}
				if len(operators) == 0 {
					return nil, fmt.Errorf("wrong((")
				}
				operators = operators[:len(operators)-1]
			}
		}
	}
	if num.Len() > 0 {
		result = append(result, num.String())
	}
	for len(operators) > 0 {
		result = append(result, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}
	fmt.Println(result)
	return result, nil
}

func evaluate(parsedExpression []string, parentId int, repo *repo.TasksRepository) (float64, error) {
	var stack []float64

	for _, token := range parsedExpression {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else if isOperator2(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("not enough operands for '%s'", token)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			res, err := EvaluateSimpleExpression(a, b, token, parentId, repo)
			if err != nil {
				return 0, err
			}

			stack = append(stack, res)
		} else {
			return 0, fmt.Errorf("invalid token: %s", token)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	log.Println(stack[0])
	err := repo.UpdateExpressionResult(parentId, "successfully calculated", stack[0])
	log.Println("trying to insert result: ", stack[0])
	if err != nil {
		logger.L().Logf(0, "can't update result | err: %v", err)
	}
	return stack[0], nil
}

func isOperator2(token string) bool {
	switch token {
	case "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func Calc(expression string, resChan chan float64, errChan chan error, PID int, repo *repo.TasksRepository) {
	parsed, err := Parse(expression)
	if err != nil {
		errChan <- err
		return
	}

	res, err := evaluate(parsed, PID, repo)
	if err != nil {
		errChan <- err
		return
	}

	resChan <- res
}
