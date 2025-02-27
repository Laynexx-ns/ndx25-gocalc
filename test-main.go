package main

import (
	"finalTaskLMS/orchestrator/pkg/calc"
	"fmt"
)

func main() {
	c := make(chan float64, 1)
	ec := make(chan error, 1)

	res, _ := calc.Calc("33+128+123/5", c, ec, 3)

	fmt.Println(res)
}
