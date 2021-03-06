package assembly

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

type relInfo struct {
	relation, cardinality string
	kind, otherKind       hierarchy
}

func (p *relInfo) flush(store *relStore) (err error) {
	if len(p.relation) > 0 {
		store.list = append(store.list, *p)
	}
	return
}

// we cant read and write to the database simultaneously with a single db? object
// so we collect the desired output and write it in a loop
type relStore struct {
	list []relInfo
}

func (store *relStore) addTerm(m *Assembler) (err error) {
	for _, p := range store.list {
		if !p.kind.valid || !p.otherKind.valid {
			e := errutil.New("couldnt determine valid lowest common relation")
			err = errutil.Append(err, e)
		} else if e := m.WriteRelation(p.relation, p.kind.name,
			p.cardinality, p.otherKind.name); e != nil {
			err = errutil.Append(err, e)
		} else if e := m.WriteVerb(p.relation, ""); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// in, eph_relation: R, K, cardinality, Q
// out, mdl_rel: R, K(lca), Q(lca), cardinality
// fix? right now the coalesce allows missing kinds through,
// the behavior otherwise is Scan error on column index 5, and not particularly helpful
func AssembleRelations(asm *Assembler) (err error) {
	var store relStore
	var curr, last relInfo
	// we select by R, sorted by R, C, K, Q
	// when C differs, we error.
	// when K differs, we Lca K.
	// when Q differs, we Lca Q.
	if e := tables.QueryAll(asm.cache.DB(),
		`select
			nr.name,
			er.cardinality,
			nk.name, coalesce(ak.path, ""),
			nq.name, coalesce(aq.path, "")
		from eph_relation er 
		join eph_named nr
			on (er.idNamedRelation = nr.rowid)
		left join eph_named nk
			on (er.idNamedKind = nk.rowid)
		left join eph_named nq
			on (er.idNamedOtherkind = nq.rowid)
		left join mdl_kind ak
			on (ak.kind = nk.name)
		left join mdl_kind aq
			on (aq.kind = nq.name)
		order by nr.name, er.cardinality, nk.name, nq.name
		`, func() (err error) {
			// when R differs, write to the output.
			if last.relation != curr.relation {
				last.flush(&store)
				// move curr into last for the next queried row.
				curr.kind.set(curr.kind.getAncestry())
				curr.otherKind.set(curr.otherKind.getAncestry())
				last = curr
			} else if last.cardinality != curr.cardinality {
				// same relation can't have different cardinality(s)
				err = errutil.New("cardinality mismatch", curr.relation, last.cardinality, curr.cardinality)
			} else {
				last.kind.update(&curr.kind)
				last.otherKind.update(&curr.otherKind)
			}
			return
		},
		&curr.relation,
		&curr.cardinality,
		&curr.kind.name, &curr.kind.parents,
		&curr.otherKind.name, &curr.otherKind.parents,
	); e != nil {
		err = e
	} else {
		last.flush(&store)
		err = store.addTerm(asm)
	}
	return
}
