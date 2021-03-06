package assembly

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/lang"
)

// goal:
// . expanded hierarchy stored per kind ( kind | comma-separated-ancestors )
// considerations:
// . contradictory ancestors and cycles ( Ts->Ks->Ts )
// . kinds without a defined hierarchy ( ie. named, but not in kinds table )
// . singular kinds ( the definition of kinds should always use plural names )
// . kinds containing punctuation ( especially "," since that used for the expanded hierarchy )
// . misspellings, near spellings
func AssembleAncestry(asm *Assembler, baseKind string) (err error) {
	var kinds cachedKinds // collect all kinds
	if e := kinds.AddDescendentsOf(asm.cache.DB(), baseKind); e != nil {
		err = errutil.New("couldn't determine ancestry of", baseKind)
	} else if e := kinds.AddDescendentsOf(asm.cache.DB(), "records"); e != nil {
		err = errutil.New("couldn't determine ancestry of", "records")
	} else {
		// write ancestors
		for k, v := range kinds.cache {
			// validate k
			if strings.ToLower(k) != k && len(k) != 2 {
				e := errutil.Fmt("BUG: there are issues with mixed case kinds (ex. %q)", k)
				err = errutil.Append(err, e)
			} else if lang.HasBadPunct(k) {
				e := errutil.New("kind shouldn't contain punctuation", k)
				err = errutil.Append(err, e)
			} else if e := asm.WriteAncestor(k, v.GetAncestors()); e != nil {
				// fix? do we want to store kinds as all "uppercase"
				// fix? future? mispellings? ( or leave that to a spellcheck in the html doc )
				err = errutil.Append(err, e)
			}
		}
		if err == nil {
			if e := reportMissingKinds(asm); e != nil {
				err = e
			} else if asm.IssueCount > 0 {
				err = errutil.Fmt("Assembly has %d outstanding issues", asm.IssueCount)
			}
		}
	}
	return
}
