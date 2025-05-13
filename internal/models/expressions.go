package models

import "github.com/google/uuid"

type Expressions struct {
	Id         int       `json:"id"`
	UserId     uuid.UUID `json:"user_id"`
	Status     string    `json:"status"`
	Result     float64   `json:"result"`
	Expression string    `json:"expression"`
}

type UserExpressions struct {
	Expression string
}
