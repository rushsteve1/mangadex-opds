package shared

import (
	"cmp"
	"testing"
)

func Tr(m map[string]string) string {
	// ja-ro is a special locale that MangaDex uses which should always exist
	return cmp.Or(m[GlobalOptions.Language], m["ja-ro"], "Unknown")
}

func Second[T any](ar []T) (out T) {
	if len(ar) == 1 {
		return out
	}

	return ar[1]
}

func AssertEq[T comparable](t *testing.T, actual, expected T) bool {
	if actual != expected {
		t.Errorf(
			"assertion failed\n%+v\nis not equal to the expected\n%+v",
			actual,
			expected,
		)
		return false
	}
	return true
}

func AssertNeq[T comparable](t *testing.T, actual, unexpected T) bool {
	if actual == unexpected {
		t.Errorf(
			"assertion failed\n%+v\nis equal to the unexpected\n%+v",
			actual,
			unexpected,
		)
		return false
	}
	return true
}
