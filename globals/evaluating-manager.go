package globals

type ExpressionTracker struct {
	ExId             int
	Expression       string
	PrimeEvaluations []PrimeEvaluation
}

type PrimeEvaluation struct {
	ParentId  int     `json:"parentId"`
	Id        int     `json:"id"`
	Arg1      float64 `json:"arg1"`
	Arg2      float64 `json:"arg2"`
	Operation rune    `json:"operation"`
}

type PrimeEvaluationResponse struct {
	ParentID      int     `json:"parentId"`
	Id            int     `json:"id"`
	Arg1          int     `json:"arg1"`
	Arg2          int     `json:"arg2"`
	Operation     rune    `json:"operation"`
	OperationTime int     `json:"operation_time"`
	Result        float64 `json:"result"`
}
