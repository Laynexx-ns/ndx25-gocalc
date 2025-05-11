package models

type Expressions struct {
	Id         int
	Status     string
	Result     float64
	Expression string
}

type UserExpressions struct {
	Expression string
}

type ExpressionsResponse struct {
	Id     int
	Status string
	Result float64
}
