package assembly

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	r "reflect"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
	"github.com/reiver/go-porterstemmer"
)

type IssueReport func(pos reader.Position, msg string)

func cat(str ...string) string {
	return strings.Join(str, " ")
}

func NewAssembler(db *sql.DB) *Assembler {
	reportNothing := func(reader.Position, string) {}
	return NewAssemblerReporter(db, reportNothing)
}

func NewAssemblerReporter(db *sql.DB, report IssueReport) *Assembler {
	return &Assembler{tables.NewCache(db), report, 0}
}

type Assembler struct {
	cache      *tables.Cache
	issueFn    IssueReport
	IssueCount int
}

func (m *Assembler) reportIssue(src, ofs, msg string) {
	pos := reader.Position{Source: src, Offset: ofs}
	m.issueFn(pos, msg)
	m.IssueCount++
}
func (m *Assembler) reportIssuef(src, ofs, fmt string, args ...interface{}) {
	m.reportIssue(src, ofs, errutil.Sprintf(fmt, args...))
}

// write kind and comma separated ancestors
func (m *Assembler) WriteAncestor(kind, path string) (err error) {
	_, e := m.cache.Exec(mdl_kind, kind, path)
	return e
}

func (m *Assembler) WriteCheck(name, testType, expect string) error {
	_, e := m.cache.Exec(mdl_check, name, testType, expect)
	return e
}

func (m *Assembler) WriteField(kind, field, fieldType string) error {
	_, e := m.cache.Exec(mdl_field, kind, field, fieldType)
	return e
}

// WriteDefault: if no specific value has been assigned to the an instance of the idModelField's kind,
// the passed default value will be used for that instance's kind.
func (m *Assembler) WriteDefault(kind, field string, value interface{}) error {
	_, e := m.cache.Exec(mdl_default, kind, field, value)
	return e
}

func DomainNameOf(domain, noun string) string {
	var b strings.Builder
	b.WriteRune('#')
	if len(domain) > 0 {
		b.WriteString(domain)
		b.WriteString("::")
	}
	b.WriteString(lang.Camelize(noun))
	return b.String()
}

func (m *Assembler) WriteNoun(noun, kind string) error {
	_, e := m.cache.Exec(mdl_noun, noun, kind)
	return e
}

// WriteName for noun
func (m *Assembler) WriteName(noun, name string, rank int) error {
	_, e := m.cache.Exec(mdl_name, noun, name, rank)
	return e
}

// WriteNounWithNames writes the noun to the model,
// and splits the name into separate space separated words.
// each word is recorded as a possible reference to the noun.
func (m *Assembler) WriteNounWithNames(domain, noun, kind string) (err error) {
	id := DomainNameOf(domain, noun)
	if e := m.WriteNoun(id, kind); e != nil {
		err = errutil.Append(err, e)
	} else {
		lower := strings.ToLower(noun)
		if e := m.WriteName(id, lower, 0); e != nil {
			err = errutil.Append(err, e)
		} else {
			var ofs int
			camel := lang.Camelize(noun)
			if camel != lower {
				if e := m.WriteName(id, camel, 1); e != nil {
					err = errutil.Append(err, e)
				}
				ofs++
			}
			split := lang.Fields(noun)
			if cnt := len(split); cnt > 1 {
				cnt += ofs
				for i, k := range split {
					rank := cnt - i
					if e := m.WriteName(id, strings.ToLower(k), rank); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
	}
	return
}

// see copyPatterns, mostly this isnt used.
func (m *Assembler) WritePat(name, paramName, paramType string, idx int64) error {
	_, e := m.cache.Exec(mdl_pat, name, paramName, paramType, idx)
	return e
}

func (m *Assembler) WritePlural(one, many string) error {
	_, e := m.cache.Exec(mdl_plural, one, many)
	return e
}

func (m *Assembler) WriteProg(progName, typeName string, bytes []byte) (int64, error) {
	return m.cache.Exec(mdl_prog, progName, typeName, bytes)
}

func (m *Assembler) WriteGob(progName string, cmd interface{}) (ret int64, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	rval := r.ValueOf(cmd)
	if e := enc.EncodeValue(rval); e != nil {
		err = errutil.New("WriteGob, error encoding value", e)
	} else {
		typeName := rval.Elem().Type().Name()
		ret, err = m.WriteProg(progName, typeName, buf.Bytes())
	}
	return
}

func (m *Assembler) WriteGobs(gobs map[string]interface{}) (err error) {
	for k, v := range gobs {
		if _, e := m.WriteGob(k, v); e != nil {
			err = e
			break
		}
	}
	return
}

func (m *Assembler) WriteRelation(relation, kind, cardinality, otherKind string) error {
	_, e := m.cache.Exec(mdl_rel, relation, kind, cardinality, otherKind)
	return e
}

// WriteStart: store the initial value of an instance's field used at start of play.
func (m *Assembler) WriteStart(noun, field string, value interface{}) error {
	_, e := m.cache.Exec(mdl_start, noun, field, value)
	return e
}

func (m *Assembler) WriteTrait(aspect, trait string, rank int) error {
	_, e := m.cache.Exec(mdl_aspect, aspect, trait, rank)
	return e
}

func (m *Assembler) WriteVerb(relation, verb string) error {
	const asm_verb = `insert into asm_verb(relation, stem)
				select ?1, ?2
				where not exists (
					select 1 from asm_verb v
					where v.relation=?1 and v.stem=?2
				)`
	stem := porterstemmer.StemString(verb)
	_, e := m.cache.Exec(asm_verb, relation, stem)
	return e
}

var mdl_aspect = tables.Insert("mdl_aspect", "aspect", "trait", "rank")
var mdl_check = tables.Insert("mdl_check", "name", "type", "expect")
var mdl_default = tables.Insert("mdl_default", "kind", "field", "value")
var mdl_field = tables.Insert("mdl_field", "kind", "field", "type")
var mdl_kind = tables.Insert("mdl_kind", "kind", "path")
var mdl_name = tables.Insert("mdl_name", "noun", "name", "rank")
var mdl_noun = tables.Insert("mdl_noun", "noun", "kind")
var mdl_pair = tables.Insert("mdl_pair", "noun", "relation", "otherNoun")
var mdl_pat = tables.Insert("mdl_pat", "pattern", "param", "type", "idx")
var mdl_plural = tables.Insert("mdl_plural", "one", "many")
var mdl_prog = tables.Insert("mdl_prog", "name", "type", "bytes")
var mdl_rel = tables.Insert("mdl_rel", "relation", "kind", "cardinality", "otherKind")
var mdl_spec = tables.Insert("mdl_spec", "type", "name", "spec")
var mdl_start = tables.Insert("mdl_start", "noun", "field", "value")
