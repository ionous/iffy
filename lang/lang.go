package lang

import (
	"bitbucket.org/pkg/inflect"
	"regexp"
	"strings"
	"unicode"
)

var Articles = []string{"the", "a", "an", "our", "some"}
var articleBar = strings.Join(Articles, "|")
var articles = regexp.MustCompile(`^((?i)` + articleBar + `)\s`)
var articleBare = regexp.MustCompile("^(" + articleBar + ")$")

const NewLine = "\n"
const Space = " "

func IsArticle(s string) bool {
	return articleBare.MatchString(s)
}

func SliceArticle(str string) (article, bare string) {
	n := strings.TrimSpace(str)
	if pair := articles.FindStringIndex(n); pair == nil {
		bare = n
	} else {
		split := pair[1] - 1
		article = n[:split]
		bare = strings.TrimSpace(n[split:])
	}
	return article, bare
}

func StripArticle(str string) string {
	_, bare := SliceArticle(str)
	return bare
}

//
func Singularize(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Singularize(s)
	}
	return
}

//
func Pluralize(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Pluralize(s)
	}
	return
}

// Capitalize returns a new string, starting the first word with a capital.
func Capitalize(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Capitalize(s)
	}
	return
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

func Lowerspace(s string) (ret string) {
	if len(s) > 0 {
		res := inflect.Humanize(s)
		ret = Lowercase(res)
	}
	return
}

func Lowercase(s string) string {
	return strings.ToLower(s)
}

func StartsWith(s string, set ...string) (ok bool) {
	for _, x := range set {
		if strings.HasPrefix(s, x) {
			ok = true
			break
		}
	}
	return ok
}

//http://www.mudconnect.com/SMF/index.php?topic=74725.0
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
