package types

import "testing"

func TestNonEmptyBanner(t *testing.T) {
	if Banner == "" || len(Banner) == 0 {
		t.Fatal("Empty Banner constant")
	}
}
