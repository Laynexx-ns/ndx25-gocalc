package calc

import (
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

var o *types.Orchestrator
var resultChans = make(map[int]chan float64)
var parentId int

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/'
}

func AddToQueue(a, b float64, operator string, c chan float64, id int, o *types.Orchestrator) {
	fmt.Println("invoked add to queue")
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if _, ok := o.Subs[parentId]; !ok {
		o.Subs[parentId] = make(chan struct{}, 1)
	}

	resultChans[id] = make(chan float64, 1)

	o.Queue = append(o.Queue, globals.PrimeEvaluation{
		ParentID:      parentId,
		Id:            id,
		Arg1:          a,
		Arg2:          b,
		Operation:     operator,
		OperationTime: 0,
		Result:        0,
	})
	fmt.Printf("Task added to queue: %+v\n", globals.PrimeEvaluation{
		ParentID:  parentId,
		Id:        id,
		Arg1:      a,
		Arg2:      b,
		Operation: operator,
	})

	go WatchQueue(parentId, c, id, o)
}

func WatchQueue(parentId int, c chan float64, id int, o *types.Orchestrator) {
	sub := o.Subs[parentId]
	for {
		select {
		case <-sub:
			o.Mu.Lock()
			for _, v := range o.Queue {
				if v.ParentID == parentId && v.Id == id && v.OperationTime != 0 {
					c <- v.Result
					close(resultChans[id])
					delete(resultChans, id)
					o.Mu.Unlock()
					return
				}
			}
			o.Mu.Unlock()
		case <-time.After(5 * time.Second):
			o.Mu.Lock()
			close(c)
			delete(o.Subs, parentId)
			o.Mu.Unlock()
			fmt.Println("timeout reached (5s)")
			return
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

func evaluate(parsedExpression []string) (float64, error) {
	stack := []float64{}
	//ch := make(chan float64)
	id := 0

	for _, token := range parsedExpression {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else if isOperator2(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("недостаточно операндов для оператора %s", token)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			resultChan := make(chan float64, 1)
			AddToQueue(a, b, token, resultChan, id, o)
			id++

			select {
			case res := <-resultChan:
				stack = append(stack, res)
			case <-time.After(5 * time.Second):
				return 0, fmt.Errorf("таймаут операции %s", token)
			}
		} else {
			return 0, fmt.Errorf("неверный токен: %s", token)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("неверное выражение")
	}

	return stack[0], nil
}

func isOperator2(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func Calc(expression string, resChan chan float64, errChan chan error, PID int, orch *types.Orchestrator) {
	fmt.Println("calculation invoked")
	o = orch
	parentId = PID
	a, err := Parse(expression)
	if err != nil {
		errChan <- err

	}
	res, err := evaluate(a)
	if err != nil {
		errChan <- err

	}
	if resChan != nil {
		resChan <- res
	}

}
