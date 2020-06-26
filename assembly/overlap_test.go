package assembly

import (
	"reflect"
	"testing"
)

// Test the lowest common ancestor helper function: findOverlap()
func TestOverlap(t *testing.T) {
	// lca, no match
	if cmp, over := findOverlap(
		[]string{"a", "b", "c", "d", "e"},
		[]string{"z", "y", "x", "w", "e"}); cmp != 0 {
		t.Fatal("expected no match", cmp)
	} else if expected := []string{"e"}; !reflect.DeepEqual(over, expected) {
		t.Fatal("want:", expected, "have:", over)
	}
	// left match
	if cmp, over := findOverlap(
		[]string{"d", "e"},
		[]string{"c", "d", "e"}); cmp != -1 {
		t.Fatal("expected left match", cmp)
	} else if expected := []string{"d", "e"}; !reflect.DeepEqual(over, expected) {
		t.Fatal("want:", expected, "have:", over)
	}
	// right same
	if cmp, over := findOverlap(
		[]string{"c", "d", "e"},
		[]string{"d", "e"}); cmp != 1 {
		t.Fatal("expected right match", cmp)
	} else if expected := []string{"d", "e"}; !reflect.DeepEqual(over, expected) {
		t.Fatal("want:", expected, "have:", over)
	}
	// no overlap, same lengths
	if cmp, over := findOverlap(
		[]string{"x", "y", "z"},
		[]string{"a", "d", "e"}); cmp != 0 {
		t.Fatal("expected no match", cmp)
	} else if len(over) != 0 {
		t.Fatal("got:", over)
	}
	// no overlap, differing lengths
	if cmp, over := findOverlap(
		[]string{"y", "z"},
		[]string{"a", "d", "e"}); cmp != 0 {
		t.Fatal("expected no match", cmp)
	} else if len(over) != 0 {
		t.Fatal("got:", over)
	}
	// both the same
	if cmp, over := findOverlap(
		[]string{"a", "b"},
		[]string{"a", "b"}); cmp != 1 {
		t.Fatal("expected match", cmp)
	} else if len(over) != 2 {
		t.Fatal("got:", over)
	}
}

func TestOverlapLca(t *testing.T) {
	match := func(a, b, c []string) bool {
		_, chain := findOverlap(a, b)
		return reflect.DeepEqual(chain, c)
	}
	if !match([]string{"A"}, []string{"A"}, []string{"A"}) {
		t.Fatal("expected lowest common ancestor A")
	} else if !match([]string{"A"}, []string{"B", "A"}, []string{"A"}) {
		t.Fatal("expected lowest common ancestor A")
	} else if !match([]string{"B", "A"}, []string{"B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor A")
	} else if !match([]string{"D", "C", "B", "A"}, []string{"B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor B")
	} else if !match([]string{"B", "A"}, []string{"D", "C", "B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor B")
	} else if !match([]string{"E", "F", "B", "A"}, []string{"D", "C", "B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor B")
	} else if !match([]string{"D", "C", "B", "A"}, []string{"E", "F", "B", "A"}, []string{"B", "A"}) {
		t.Fatal("expected lowest common ancestor B")
	} else if !match([]string{"D", "P", "T"}, []string{"T"}, []string{"T"}) {
		t.Fatal("expected lowest common ancestor T")
	} else if !match([]string{"D", "E", "F"}, []string{"C", "B", "A"}, nil) {
		t.Fatal("expected no lowest common ancestor")
	}
}
