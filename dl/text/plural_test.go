package text

import (
	"testing"
)

func TestPlurals(t *testing.T) {
	test := map[string]string{
		"apple": "apples",
		"movie": "movies",
		"man":   "men",
		"ox":    "oxen",
		"purse": "purses",
		"rice":  "rice",
	}
	p := make(Plurals)
	//
	try := func(single, expected string) {
		got := p.Pluralize(single)
		if got != expected {
			t.Fatal(single, "expected", expected, "got", got)
		}
	}
	for single, expected := range test {
		try(single, expected)
	}
	single, expected := "purse", "clutch"
	p.AddPlural(single, expected)
	try(single, expected)
}
