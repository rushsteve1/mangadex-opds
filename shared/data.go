package shared

import "cmp"

type Data[T any] struct {
	Result string `json:"result"`
	Data   T      `json:"data"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func Tr(m map[string]string) string {
	lang := cmp.Or(GlobalOptions.Language, "en")
	return cmp.Or(m[lang], m["en"], m["en-ro"], "Unknown")
}
