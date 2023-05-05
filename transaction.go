package models

import "time"

type Transaction struct {
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type Statistics struct {
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
	Count int64   `json:"count"`
}

type Location struct {
	City string `json:"city"`
}
