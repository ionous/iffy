package lang

import (
	"testing"
)

func TestStripArticle(t *testing.T) {
	type Pair struct {
		src, article, text string
	}
	p := []Pair{
		{src: "the evil fish", article: "the", text: "evil fish"},
		{src: "The Capital", article: "The", text: "Capital"},
		{src: "some fish", article: "some", text: "fish"},
		{src: " a   space ", article: "a", text: "space"},
		{src: "dune, a desert planet", article: "", text: "dune, a desert planet"},
	}

	for _, p := range p {
		a, b := SliceArticle(p.src)
		if p.text != b {
			t.Fatalf("text: '%s'", p.src)
		}
		if p.article != a {
			t.Fatalf("text: '%s'", p.src)
		}
	}

	if "article" != StripArticle("the article") {
		t.Fatal("article")
	}
}

// ensure package.inflect works as expected for a few cases....
// fix: they dont really work they way id expect.
func TestInflect(t *testing.T) {
	p := []struct {
		test, want string
		format     func(s string) string
	}{
		{
			"boop",
			"Boop",
			Capitalize,
		},
		{
			"BOOP",
			"Boop",
			Capitalize,
		}, {
			"another day. at SEA.... oh my.",
			"Another Day. At S E A.... Oh My.",
			Titlecase,
		}, {
			"another day. at SEA.... oh my.",
			"Another day. At sea.... Oh my.",
			SentenceCase,
		},
	}
	for i, p := range p {
		got := p.format(p.test)
		if got != p.want {
			t.Fatalf("test %v failed, got %q", i, got)
		}
	}
}
