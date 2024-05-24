package models

// Data is a wrapping type that is used to parse most responses from the MangaDex API
type Data[T any] struct {
	Result string `json:"result"`
	Data   T      `json:"data"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
