package assembly

import (
	"reflect"
	"testing"
)

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
