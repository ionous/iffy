package lang

import (
	"testing"
)

// go test --run TestVowels
func TestVowels(t *testing.T) {
	if !StartsWithVowel("evil fish") {
		t.Fatal("error")
	}
}
