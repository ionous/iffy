package chart

import (
	"strings"
	"testing"
)

func TestSpaces(t *testing.T) {
	count := func(x int, str string) {
		// the fail point is one-indexed,
		n, _ := innerParse(spaces, str)
		if n-1 != x {
			t.Fatal(str, "len:", n-1)
		}
	}
	count(0, "a")
	count(-1, "")
	count(-1, strings.Repeat(" ", 5))
	count(3, strings.Repeat(" ", 3)+"x")
}
