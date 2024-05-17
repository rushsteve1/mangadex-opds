package models

import (
	"cmp"

	"github.com/rushsteve1/mangadex-opds/shared"
)

type Data[T any] struct {
	Result string `json:"result"`
	Data   T      `json:"data"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func Tr(m map[string]string) string {
	// ja-ro is a special locale that MangaDex uses which should always exist
	return cmp.Or(m[shared.GlobalOptions.Language], m["ja-ro"], "Unknown")
}
