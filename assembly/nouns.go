package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
)

// output:
// - mdl_noun: noun(int), kind, [ scene ]
// - mdl_name: noun, name part, rank
//
// inputs:
// - mdl_kind: kind, path
// - eph_noun: noun, kind
// - eph_name: for nouns.
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

func (store *nounStore) write(m *Modeler) (err error) {
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

func DetermineNouns(m *Modeler, db *sql.DB) (err error) {
	var store nounStore
	var curr, last nounInfo
	if e := dbutil.QueryAll(db,
		`select nn.name, nk.name, coalesce(ak.path, "")
		from eph_noun n join eph_named nn
			on (n.idNamedNoun = nn.rowid)
		left join eph_named nk
			on (n.idNamedKind = nk.rowid)
		left join mdl_kind ak
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
		err = store.write(m)
	}
	return
}
