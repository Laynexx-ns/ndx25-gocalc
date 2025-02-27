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

func AddToQueue(a, b float64, operator rune, c chan float64, id int, o *types.Orchestrator) {
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

func NotifySubscribers(parentId int) {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if sub, ok := o.Subs[parentId]; ok {
		close(sub)
		delete(o.Subs, parentId)
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
			AddToQueue(a, b, rune(v[0]), ch, id, o)
			select {
			case res := <-ch:
				evex = append(evex[:i-2], append([]string{fmt.Sprintf("%f", res)}, evex[i+1:]...)...)
				i -= 2
			case <-time.After(5 * time.Second):
				return 0, fmt.Errorf("таймаут агента")
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
	resChan <- res

}
