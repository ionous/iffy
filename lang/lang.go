package lang

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/ionous/inflect"
)

var Articles = []string{"the", "a", "an", "our", "some"}
var articleBar = strings.Join(Articles, "|")
var articles = regexp.MustCompile(`^((?i)` + articleBar + `)\s`)
var articleBare = regexp.MustCompile("^(" + articleBar + ")$")

const NewLine = "\n"
const Space = " "

// IsArticle returns true if the passed string starts with one of the common determiners.
func IsArticle(s string) bool {
	return articleBare.MatchString(s)
}

// SliceArticle splits the passed string into a common determiner (if any) and the remaining text.
func SliceArticle(s string) (article, bare string) {
	n := strings.TrimSpace(s)
	if pair := articles.FindStringIndex(n); pair == nil {
		bare = n
	} else {
		split := pair[1] - 1
		article = n[:split]
		bare = strings.TrimSpace(n[split:])
	}
	return article, bare
}

// StripArticle removes common determiners from the start of the passed string.
func StripArticle(s string) string {
	_, bare := SliceArticle(s)
	return bare
}

// Singularize attempts to return the singular form of the passed assumed plural string.
func Singularize(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Singularize(s)
	}
	return
}

// Pluralize attempts to return the plural form of the passed assumed singular string.
func Pluralize(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Pluralize(s)
	}
	return
}

// IsPlural returns true if the passed string seems pluralized.
func IsPlural(s string) bool {
	return s != inflect.Singularize(s)
}

// Capitalize returns a new string, starting the first word with a capital.
func Capitalize(s string) (ret string) {
	if len(s) > 0 {
		// fix? capitalize doesnt handle leading spaces well.
		// what should it do?
		var lead string
		if i := strings.IndexFunc(s, func(u rune) bool {
			return !unicode.IsSpace(u)
		}); i >= 0 {
			lead, s = s[:i], s[i:]
		}
		ret = lead + inflect.Capitalize(strings.ToLower(s))
	}
	return
}

// SentenceCase returns the passed string in lowercase, starting new sentences with capital letters.
func SentenceCase(s string) string {
	sentences := strings.Split(s, ". ")
	for i, s := range sentences {
		sentences[i] = Capitalize(s)
	}
	return strings.Join(sentences, ". ")
}

// IsCapitalized returns true if the passed string starts with an upper case letter
func IsCapitalized(n string) (ret bool) {
	for _, r := range n {
		ret = unicode.IsUpper(r)
		break
	}
	return
}

// Titlecase returns a new string, starting every word with a capital.
func Titlecase(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Titleize(s)
	}
	return
}

// Lowerspace returns the passed string in lowercase with common word separators changed into spaces.
func Lowerspace(s string) (ret string) {
	if len(s) > 0 {
		res := inflect.Humanize(s)
		ret = Lowercase(res)
	}
	return
}

// Lowercase is an alias for strings.ToLower
func Lowercase(s string) string {
	return strings.ToLower(s)
}

// StartsWith returns true if the passed string starts with any one of the passed strings in set
func StartsWith(s string, set ...string) (ok bool) {
	for _, x := range set {
		if strings.HasPrefix(s, x) {
			ok = true
			break
		}
	}
	return ok
}

// StartsWithVowel returns true if the passed strings starts with a vowel or vowel sound.
// http://www.mudconnect.com/SMF/index.php?topic=74725.0
func StartsWithVowel(str string) (vowelSound bool) {
	s := strings.ToUpper(str)
	if StartsWith(s, "A", "E", "I", "O", "U") {
		if !StartsWith(s, "EU", "EW", "ONCE", "ONE", "OUI", "UBI", "UGAND", "UKRAIN", "UKULELE", "ULYSS", "UNA", "UNESCO", "UNI", "UNUM", "URA", "URE", "URI", "URO", "URU", "USA", "USE", "USI", "USU", "UTA", "UTE", "UTI", "UTO") {
			vowelSound = true
		}
	} else if StartsWith(s, "HEIR", "HERB", "HOMAGE", "HONEST", "HONOR", "HONOUR", "HORS", "HOUR") {
		vowelSound = true
	}
	return vowelSound
}

// ContainsPunct returns true if any rune of s returns true for unicode.IsPunct.
func ContainsPunct(s string) bool {
	return strings.IndexFunc(s, unicode.IsPunct) >= 0
}
