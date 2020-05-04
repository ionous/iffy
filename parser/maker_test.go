package parser_test

import (
	"strings"
	"testing"

	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
)

func TestPhraseMaker(t *testing.T) {
	if diff := pretty.Diff(sliceOf.String("look"), Phrases("look")); len(diff) > 0 {
		t.Fatal(diff)
	}
	if diff := pretty.Diff(sliceOf.String("look", "l"), Phrases("look/l")); len(diff) > 0 {
		t.Fatal(diff)
	}
	if diff := pretty.Diff([]string{
		"look inside something", "look inside wicked",
		"look in something", "look in wicked",
		"l inside something", "l inside wicked",
		"l in something", "l in wicked"},
		Phrases("look/l inside/in something/wicked")); len(diff) > 0 {
		t.Fatal(diff)
	}
}

// generate permutations from inform-like slash phrases
func Phrases(phrase string) (ps []string) {
	// step 1 split the phrase into space chunks
	var multi [][]string
	for _, f := range strings.Fields(phrase) {
		multi = append(multi, strings.Split(f, "/"))
	}
	for _, m := range multi {
		ps = permute(ps, m)
	}
	return
}

func permute(a, b []string) (ps []string) {
	if len(a) == 0 {
		ps = b
	} else {
		for _, a := range a {
			for _, b := range b {
				n := strings.Join([]string{a, b}, " ")
				ps = append(ps, n)
			}
		}
	}
	return
}
