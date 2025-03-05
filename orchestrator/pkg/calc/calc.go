package calc

import (
	"errors"
	"finalTaskLMS/globals"
	"finalTaskLMS/orchestrator/types"
	"fmt"
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

func EvaluateSimpleExpression(a, b float64, operand string, parentId int, orch *types.Orchestrator) (float64, error) {
	resChan := make(chan float64, 1)
	errChan := make(chan error, 1)
	defer close(resChan)
	defer close(errChan)

	orch.Mu.Lock()
	id := len(orch.Queue)

	orch.Queue = append(orch.Queue, globals.PrimeEvaluation{
		ParentID:      parentId,
		Id:            id,
		Arg1:          a,
		Arg2:          b,
		Operation:     operand,
		OperationTime: 0,
		Result:        0,
	})
	orch.Mu.Unlock()

	WaitOperationResult(resChan, errChan, id, parentId, orch)

	select {
	case res := <-resChan:
		return res, nil
	case err := <-errChan:
		return 0, err
	}
}
func WaitOperationResult(resChan chan float64, errChan chan error, id, parentId int, orch *types.Orchestrator) {
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Print("OPERATION TIMEOUT")
			errChan <- errors.New("operation timeout")
			return
		case <-ticker.C:
			orch.Mu.Lock()
			if id < len(orch.Queue) && orch.Queue[id].OperationTime != 0 {
				res := orch.Queue[id].Result
				fmt.Print("OPERATION EVALUATED")
				orch.Mu.Unlock()
				resChan <- res
				return
			}
			orch.Mu.Unlock()
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

func evaluate(parsedExpression []string, parentId int, orch *types.Orchestrator) (float64, error) {
	var stack []float64
	var evRes float64

	for _, token := range parsedExpression {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else if isOperator2(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("%s", token)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			evRes, err = EvaluateSimpleExpression(a, b, token, parentId, orch)
			if err != nil {
				return 0, err
			}
			stack = append(stack, evRes)
		} else {
			return 0, fmt.Errorf("invalid operation: %s", token)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid operation")
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

func Calc(expression string, resChan chan float64, errChan chan error, PID int, orch *types.Orchestrator) {
	fmt.Println("calculation invoked")

	a, err := Parse(expression)
	if err != nil {
		errChan <- err

	}
	res, err := evaluate(a, PID, orch)
	if err != nil {
		errChan <- err

	}
	if resChan != nil {
		resChan <- res
	}

}
