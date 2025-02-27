package calculator

import (
	"os"
	"strconv"
	"time"
)

func Evaluate(arg1, arg2 float64, operand rune) float64 {
	var res float64
	duration, err := strconv.Atoi(os.Getenv("qwe"))
	if err != nil {
		duration = 1
	}

	switch operand {
	case '+':
		time.Sleep(time.Duration(duration) * time.Second)
		res = arg1 + arg2
	case '-':
		time.Sleep(time.Duration(duration) * time.Second)
		res = arg1 - arg2
	case '/':
		time.Sleep(time.Duration(duration) * time.Second)
		res = arg1 / arg2
	case '*':
		time.Sleep(time.Duration(duration) * time.Second)
		res = arg1 * arg2
	}
	return res
}
