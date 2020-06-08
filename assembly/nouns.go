package assembly

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

// output:
// - mdl_noun: noun(int), kind, [ scene ]
// - mdl_name: noun, name part, rank
//
// inputs:
// - mdl_kind: kind, path
// - eph_noun: noun, kind
// - eph_name: for nouns.
//
func AssembleNouns(asm *Assembler) (err error) {
	var store nounStore
	var curr, last nounInfo
	if e := tables.QueryAll(asm.cache.DB(),
		// note: nk is known to refer to kinds b/c it comes from eph_noun.idNamedKind
		// therefore, we dont have to filter where category=kind(s).
		`select nn.name, nk.name, coalesce(ak.path, "")
		from eph_noun en 
		join eph_named nn
			on (en.idNamedNoun = nn.rowid)
		left join eph_named nk
			on (en.idNamedKind = nk.rowid)
		join mdl_kind ak
			on (ak.kind = nk.name)
		order by nn.name, nk.name
		`, func() (err error) {
			// when the noun differs, write to the output.
			if last.noun != curr.noun {
				last.flush(&store)
				// move curr into last for the next queried row.
				curr.kind.set(curr.kind.getAncestry())
				last = curr
			} else {
				last.kind.update(&curr.kind)
			}
			return
		},
		&curr.noun, &curr.kind.name, &curr.kind.parents,
	); e != nil {
		err = e
	} else {
		last.flush(&store)
		if e := store.writeNouns(asm); e != nil {
			err = e
		} else if e := reportMissingNouns(asm); e != nil {
			err = e
		} else if asm.IssueCount > 0 {
			err = errutil.Fmt("Assembly has %d outstanding issues", asm.IssueCount)
		}
	}
	return
}

type nounInfo struct {
	noun string
	kind hierarchy
}

func (p *nounInfo) flush(store *nounStore) {
	if len(p.noun) > 0 {
		store.list = append(store.list, *p)
	}
}

// we cant read and write to the database simultaneously with a single db? object
// so we collect the desired output and write it in a loop
type nounStore struct {
	list []nounInfo
}

func (store *nounStore) writeNouns(m *Assembler) (err error) {
	for _, p := range store.list {
		if !p.kind.valid {
			e := errutil.New("couldnt determine valid lowest common ancestor")
			err = errutil.Append(err, e)
		} else if e := m.WriteNounWithNames(p.noun, p.kind.name); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
