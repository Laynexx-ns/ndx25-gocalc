package models

// ExerciseProcessingStatus / - enum для статусов обработки упражнений

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
