package globals

type ExpressionTracker struct {
	ExId             int
	Expression       string
	PrimeEvaluations []PrimeEvaluation
}

type PrimeEvaluation struct {
	ParentID      int     `json:"parentId"`
	Id            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
	Result        float64 `json:"result"`
}
