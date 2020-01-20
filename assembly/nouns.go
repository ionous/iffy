package assembly

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
)

// output:
// - mdl_noun: id, kind, [ scene ]
// - mdl_name: inst.id, name part, rank
//
// inputs:
// - mdl_kind: kind, path
// - eph_noun: id, kind
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

func (store *nounStore) write(w *Modeler) (err error) {
	for _, p := range store.list {
		if !p.kind.valid {
			e := errutil.New("couldnt determine valid lowest common ancestor")
			err = errutil.Append(err, e)
		} else if n, e := w.WriteNoun(p.kind.name); e != nil {
			err = errutil.Append(err, e)
		} else {
			w.WriteName(n, p.noun, 0)
			split := strings.Fields(p.noun)
			if cnt := len(split); cnt > 1 {
				for i, k := range split {
					rank := cnt - i
					w.WriteName(n, k, rank)
				}
			}
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
