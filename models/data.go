package models

type Data[T any] struct {
	Result string `json:"result"`
	Data   T      `json:"data"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
