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

func AddToQueue(a, b float64, operator rune, c chan float64, id int) {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if _, exists := resultChans[id]; !exists {
		resultChans[id] = make(chan float64, 1)
	}
	resultChans[id] = make(chan float64, 1)

	o.Queue = append(o.Queue, globals.PrimeEvaluation{
		ParentId:  parentId,
		Id:        id,
		Arg1:      a,
		Arg2:      b,
		Operation: operator,
	})

	go func(id int) {
		select {
		case res := <-resultChans[id]:
			c <- res
		}
	}(id)
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
	var evex []string
	ch := make(chan float64)
	id := 0
	evex = parsedExpression
	for i, v := range parsedExpression {
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			a, errA := strconv.ParseFloat(evex[i-2], 64)
			b, errB := strconv.ParseFloat(evex[i-1], 64)
			if errA != nil || errB != nil {
				return 0, fmt.Errorf("чет ты какую-то фигню на сервер отправил :(. А я, бездарь не обработал :((")
			}
			AddToQueue(a, b, rune(v[0]), ch, id)
			select {
			case res := <-ch:
				evex = append(evex[:i-2], append([]string{fmt.Sprintf("%f", res)}, evex[i+1:]...)...)
				i -= 2
			case <-time.After(5 * time.Second):
				return 0, fmt.Errorf("время ожидания агента вышло, сворачиваемся")
			}
		}
		id++

	}
	if len(evex) != 1 {
		return 0, fmt.Errorf("мой калькулятор не работает, тильт")
	}

	result, err := strconv.ParseFloat(evex[0], 64)
	if err != nil {
		return 0, fmt.Errorf("мой калькулятор не работает, тильт")
	}

	return result, nil
}

func Calc(expression string, resChan chan float64, errChan chan error, PID int, orch *types.Orchestrator) (float64, error) {
	o = orch
	parentId = PID
	a, err := Parse(expression)
	if err != nil {
		errChan <- err
		return 0, err
	}
	res, err := evaluate(a)
	if err != nil {
		errChan <- err
		return 0, err
	}
	resChan <- res
	return res, nil

}
