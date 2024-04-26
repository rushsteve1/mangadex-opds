package shared

// This is a simple type alias but it makes it clearer that
// the string should not be treated like a normal string
type UUID = string

type Data[T any] struct {
	Result string `json:"result"`
	Data   T      `json:"data"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
