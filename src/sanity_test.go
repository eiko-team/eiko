package sanity_test

import (
	"eiko/src"
	"testing"
)

func TestAbs(t *testing.T) {
	got := sanity.Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}
