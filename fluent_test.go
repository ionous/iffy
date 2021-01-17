package iffy_test

import (
	r "reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/debug"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/rel"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
	"github.com/kr/pretty"
)

// until template parsing gets re-written we cant handle fluid specs ( selector messaging )
// we can do a basic test to ensure it's possible to build the function signatures from the composer.Spec(s) tho.
func TestFluid(t *testing.T) {
	if got := makeSig((*core.IsKindOf)(nil)); !got.equals(
		"kindOf:is:",
	) {
		t.Error(got)
	}
	if got := makeSig((*core.IsExactKindOf)(nil)); !got.equals(
		"kindOf:isExactly:",
	) {
		t.Error(got)
	}
	if got := makeSig((*core.CompareNum)(nil)); !got.equals(
		"is num:equalTo:", "is num:otherThan:", "is num:greaterThan:", "is num:lessThan:", "is num:atLeast:", "is num:atMost:",
	) {
		t.Error(got)
	}
	if got := makeSig((*core.IsNotTrue)(nil)); !got.equals(
		"not:",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.Range)(nil)); !got.equals(
		"range:", "range:from:", "range:byStep:", "range:from:byStep:",
	) {
		t.Error(got)
	}
	if got := makeSig((*core.While)(nil)); !got.equals(
		"repeating while:do:",
	) {
		t.Error(got)
	}
	// rel
	if got := makeSig((*rel.RelativesOf)(nil)); !got.equals(
		"relatives:of:",
	) {
		t.Error(got)
	}
	if got := makeSig((*rel.Relate)(nil)); !got.equals(
		"relate:to:via:",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.Each)(nil)); !got.equals(
		"repeating across:asNum:do:",
		"repeating across:asTxt:do:",
		"repeating across:asRec:do:",
		"repeating across:asNum:do:elseIfEmptyDo:",
		"repeating across:asTxt:do:elseIfEmptyDo:",
		"repeating across:asRec:do:elseIfEmptyDo:",
	) {
		t.Error(got)
	}
	if got := makeSig((*core.ChooseAction)(nil)); !got.equals(
		"if:do:",
		"if:do:elseIf:", // fix: sort!?
		"if:do:elseDo:",
	) {
		t.Error(got)
	}
	if got := makeSig((*core.Assign)(nil)); !got.equals(
		"let:be:",
	) {
		t.Error(got)
	}
	if got := makeSig((*core.PutAtField)(nil)); !got.equals(
		"put:intoRec:atField:",
		"put:intoObj:atField:",
		"put:intoObjNamed:atField:",
	) {
		t.Error(got)
	}
	// list:
	if got := makeSig((*list.PutIndex)(nil)); !got.equals(
		"put:intoNumList:atIndex:",
		"put:intoRecList:atIndex:",
		"put:intoTxtList:atIndex:",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.PutEdge)(nil)); !got.equals(
		"put:intoNumList:atBack|atFront!",
		"put:intoRecList:atBack|atFront!",
		"put:intoTxtList:atBack|atFront!",
	) {
		t.Error(got)
	}
	if got := makeSig((*debug.Log)(nil)); !got.equals(
		"log:note|toDo|warning|fix!",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.Erasing)(nil)); !got.equals(
		"erasing:fromNumList:atIndex:as:do:",
		"erasing:fromRecList:atIndex:as:do:",
		"erasing:fromTxtList:atIndex:as:do:",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.EraseIndex)(nil)); !got.equals(
		"erase:fromNumList:atIndex:",
		"erase:fromRecList:atIndex:",
		"erase:fromTxtList:atIndex:",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.EraseEdge)(nil)); !got.equals(
		"erase:atBack|atFront!",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.Gather)(nil)); !got.equals(
		"gather:fromNumList:using:",
		"gather:fromRecList:using:",
		"gather:fromTxtList:using:",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.SortText)(nil)); !got.equals(
		"sort text:ascending|descending!includeCase|ignoreCase!",
		"sort text:byField:ascending|descending!includeCase|ignoreCase!",
	) {
		t.Error(got)
	}
	if got := makeSig((*list.SortRecords)(nil)); !got.equals(
		"sort records:using:",
	) {
		t.Error(got)
	}
}

func (sig signature) equals(bs ...string) bool {
	return len(pretty.Diff(sig, signature(bs))) == 0
}

func (sig signature) String() string {
	return pretty.Sprint(sig)
}

func makeSig(v composer.Composer) signature {
	sig := signature{specName(v)}
	rtype := r.TypeOf(v).Elem()
	// see also: parseSpec
	var cnt int
	export.WalkProperties(rtype, func(f *r.StructField, path []int) (done bool) {
		tags := tag.ReadTag(f.Tag)
		if _, ok := tags.Find("internal"); !ok {
			cnt++
			// write the selector:
			var label string
			var unlabeled bool
			if selector, ok := tags.Find("selector"); ok && len(selector) > 0 {
				label = selector
			} else {
				label = firstRuneLower(f.Name)
				unlabeled = ok
			}
			if cnt == 1 {
				var sep string
				if unlabeled {
					sep = ":"
				} else {
					sep = " "
				}
				sig[0] += sep
			}

			if cnt > 1 || !unlabeled {
				// very specific check for non-optional flags...
				if x, ok := r.Zero(r.PtrTo(f.Type)).Interface().(composer.Composer); ok {
					if spec := x.Compose(); spec.UsesStr() {
						if cs := spec.Strings; len(cs) > 0 {
							sig = sig.addFlags(cs)
							return // EARLY RETURN
						}
					}
				}
				switch n, k := f.Type.Name(), f.Type.Kind(); {
				default:
					//  write camel "fieldName:"
					if tags.Exists("optional") {
						sig = sig.dupSelectors(label)
					} else if unlabeled {
						// helps with "compact" to run an arg under the previous interface
						// ex. for compare num
						sig = sig.addSelector("")

					} else {
						sig = sig.addSelector(label)
					}

				case k == r.Ptr:
					// optional. so duplicate all existing selectors
					if f.Type.Elem().Kind() != r.Struct {
						sig = sig.dupSelectors(label)
						//
					} else if x, ok := r.Zero(f.Type).Interface().(composer.Composer); !ok {
						panic("pointer to struct doesnt have spec")
					} else {
						name := specName(x)
						sig = sig.dupSelectors(name)
					}

				case k == r.Interface && !strings.HasSuffix(n, "Eval") && (n != "Assignment"):
					// avoids generating (the large number of) selectors for commands matching text evals, etc.
					// assumes interfaces are all unlabeled...
					if slats := implementorsOf(f.Type); len(slats) == 0 {
						panic("no slats found") // mul will return nil
					} else {
						sig = sig.mulSelectors(slats, tags.Exists("optional"), tags.Exists("compact"))
					}
				}
			}

		}
		return
	})
	return sig
}

// return the command structs supported by the passed slot
func implementorsOf(slot r.Type) (ret []string) {
	for _, slats := range iffy.AllSlats {
		for _, slat := range slats {
			rtype := r.TypeOf(slat)
			if rtype.Implements(slot) {
				ret = append(ret, specName(slat))
			}
		}
	}
	return
}

func specName(slat composer.Composer) string {
	name := composer.SpecName(slat)
	if spec := slat.Compose(); spec.Fluent != nil {
		if len(spec.Fluent.Name) > 0 {
			name = spec.Fluent.Name
		} else {
			name = r.TypeOf(slat).Elem().Name()
		}
	}
	return camelize(name)
}

type signature []string

func (sig signature) addSelector(sel string) signature {
	for i, cnt := 0, len(sig); i < cnt; i++ {
		sig[i] = sig[i] + sel + ":"
	}
	return sig
}

// to avoid an explosion of selectors for flags, we consider them specially.
// the signature of flags may differ from how they are specified in use, tbd.
func (sig signature) addFlags(cs []string) signature {
	var b strings.Builder
	for i, c := range cs {
		writeCamel(&b, c)
		if (i + 1) < len(cs) {
			b.WriteRune('|')
		}
	}
	b.WriteRune('!')
	sel := b.String()
	for i, cnt := 0, len(sig); i < cnt; i++ {
		sig[i] = sig[i] + sel
	}
	return sig
}

func (sig signature) mulSelectors(sel []string, optional, compact bool) signature {
	var out signature
	if optional {
		out = sig
	}
	var sep string
	if !compact {
		sep = ":"
	}
	for i, cnt := 0, len(sig); i < cnt; i++ {
		for _, sel := range sel {
			out = append(out, sig[i]+sel+sep)
		}
	}
	return out
}

func (sig signature) dupSelectors(sel string) signature {
	out := sig
	for i, cnt := 0, len(sig); i < cnt; i++ {
		out = append(out, sig[i]+sel+":")
	}
	return out
}

func firstRuneLower(s string) (ret string) {
	rs := []rune(s)
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}

func camelize(s string) string {
	var b strings.Builder
	writeCamel(&b, s)
	return b.String()
}

func writeCamel(b *strings.Builder, s string) {
	nextUpper := false
	for i, u := range s {
		if i == 0 {
			b.WriteRune(unicode.ToLower(u))
		} else if u == '_' || u == ' ' {
			nextUpper = true
		} else if nextUpper {
			b.WriteRune(unicode.ToUpper(u))
			nextUpper = false
		} else {
			b.WriteRune(u)
		}
	}
}
