package parser_test

import (
	"strings"
	"testing"

	. "github.com/ionous/iffy/parser"
)

func TestWord(t *testing.T) {
	match := func(input, goal string) (Result, error) {
		match := &Word{goal}
		return match.Scan(nil, nil, Cursor{Words: strings.Fields(input)})
	}
	t.Run("match", func(t *testing.T) {
		if res, e := match("Beep", "beep"); e != nil {
			t.Fatal("error", e)
		} else if w, ok := res.(ResolvedWord); !ok {
			t.Fatalf("%T", res)
		} else if w.Word != "Beep" {
			t.Fatal(w)
		}
	})
	t.Run("mismatch", func(t *testing.T) {
		if res, e := match("Boop", "beep"); e == nil {
			t.Fatal("expected error", res)
		}
	})
}

func TestVariousOf(t *testing.T) {
	words := func(goal ...string) (ret []Scanner) {
		for _, g := range goal {
			ret = append(ret, &Word{g})
		}
		return
	}

	match := func(input string, goal Scanner) (Result, error) {
		return goal.Scan(nil, nil, Cursor{Words: strings.Fields(input)})
	}
	t.Run("any", func(t *testing.T) {
		wordList := words("beep", "blox")
		t.Run("match", func(t *testing.T) {
			if res, e := match("Beep", &AnyOf{wordList}); e != nil {
				t.Fatal("error", e)
			} else if w, ok := res.(ResolvedWord); !ok {
				t.Fatalf("%T", res)
			} else if w.Word != "Beep" {
				t.Fatal(w)
			}
			if res, e := match("Blox", &AnyOf{wordList}); e != nil {
				t.Fatal("error", e)
			} else if w, ok := res.(ResolvedWord); !ok {
				t.Fatalf("%T", res)
			} else if w.Word != "Blox" {
				t.Fatal(w)
			}
		})
		t.Run("mismatch", func(t *testing.T) {
			if res, e := match("Boop", &AnyOf{wordList}); e == nil {
				t.Fatal("expected error", res)
			}
			if res, e := match("Beep", &AnyOf{}); e == nil {
				t.Fatal("expected error", res)
			}
			if res, e := match("", &AnyOf{}); e == nil {
				t.Fatal("expected error", res)
			}
		})
	})

	t.Run("all", func(t *testing.T) {
		wordList := words("beep", "blox")
		t.Run("match", func(t *testing.T) {
			if ares, e := match("Beep BLOX", &AllOf{wordList}); e != nil {
				t.Fatal("error", e)
			} else if res, ok := ares.(*ResultList); !ok {
				t.Fatalf("%T", ares)
			} else if cnt := res.WordsMatched(); cnt != 2 {
				t.Fatal("mismatch", cnt)
			} else {
				expect := []string{"Beep", "BLOX"}
				for i, res := range res.Results() {
					if w, ok := res.(ResolvedWord); !ok {
						t.Fatalf("%T", res)
					} else if w.Word != expect[i] {
						t.Fatal(i, w.Word)
					}
				}
			}
		})
		t.Run("mismatch", func(t *testing.T) {
			if res, e := match("BLOX Beep", &AllOf{wordList}); e == nil {
				t.Fatal("expected error", res)
			}
			if res, e := match("Beep", &AllOf{wordList}); e == nil {
				t.Fatal("expected error", res)
			}
			if res, e := match("BLOX", &AllOf{wordList}); e == nil {
				t.Fatal("expected error", res)
			}
			if res, e := match("Boop", &AllOf{wordList}); e == nil {
				t.Fatal("expected error", res)
			}
			if res, e := match("", &AllOf{wordList}); e == nil {
				t.Fatal("expected error", res)
			}
			if res, e := match("", &AllOf{}); e == nil {
				t.Fatal("expected error", res)
			}
		})
	})
}
