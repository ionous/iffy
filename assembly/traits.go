package assembly

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

// goal: build table of aspects and their traits.
// considerations:
// o ambiguous ranks ( ignoring ranks for now. )
// . conflicting traits ( different aspects; future, resolve via implications? )
// . missing traits ( named but not specified )
// . missing aspects ( named but not specified )
// o misspellings, near spellings ( ex. for missing traits )
func AssembleAspects(asm *Assembler) (err error) {
	var curr, last aspectInfo
	var traits []aspectInfo // cant read and write to the db simultaneously
	if e := tables.QueryAll(asm.cache.DB(),
		`select nt.name, na.name
	from eph_trait t 
	join eph_named nt
		on (t.idNamedTrait = nt.rowid)
	left join eph_named na
		on (t.idNamedAspect = na.rowid)
	order by na.name, t.rank, nt.name`, func() (err error) {
			switch traitsMatch, aspectsMatch := last.Trait == curr.Trait, last.Aspect == curr.Aspect; {
			case traitsMatch && !aspectsMatch:
				err = errutil.New("same trait different aspect", curr.Trait, curr.Aspect, last.Aspect)

			case !traitsMatch:
				traits = append(traits, curr)
				last = curr
			}
			return
		}, &curr.Trait, &curr.Aspect); e != nil {
		err = e
	} else {
		var last string
		var rank int
		for _, t := range traits {
			if last != t.Aspect {
				last, rank = t.Aspect, 0
			} else {
				rank++
			}
			if e := asm.WriteTrait(t.Aspect, t.Trait, rank); e != nil {
				err = errutil.Append(err, e)
			}
		}
		if err == nil {
			if e := reportMissingAspects(asm); e != nil {
				err = e
			} else if asm.IssueCount > 0 {
				err = errutil.Fmt("Assembly has %d outstanding issues", asm.IssueCount)
			} else if e := reportMissingTraits(asm); e != nil {
				err = e
			} else if asm.IssueCount > 0 {
				err = errutil.Fmt("Assembly has %d outstanding issues", asm.IssueCount)
			}
		}
	}
	return
}

type aspectInfo struct {
	Aspect string
	Trait  string
}
