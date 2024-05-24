package shared

import (
	"cmp"
	"testing"
)

// Tr takes a map of translation strings and returns the one matching the language
// in [GlobalOptions] OR the special `ja-ro` translation OR the fallback string "Unknown".
func Tr(m map[string]string) string {
	// ja-ro is a special locale that MangaDex uses which should always exist
	return cmp.Or(m[GlobalOptions.Language], m["ja-ro"], "Unknown")
}

// AssertEq asserts that two values are equal in a test, logging if they are not.
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

// AssertEq asserts that two values are NOT equal in a test, logging if they are.
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
