package express

import (
	r "reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
	"github.com/ionous/iffy/lang"
	"github.com/kr/pretty"
)

// until template parsing gets re-written we cant handle fluid specs ( selector messaging )
// we can do a basic test to ensure it's possible to build the function signatures from the composer.Spec(s) tho.
func TestFluid(t *testing.T) {
	v := (*core.PutAtField)(nil)
	rtype := r.TypeOf(v).Elem()
	spec := v.Compose()
	fluid := spec.Fluent

	short := shortName(rtype, fluid.Name)

	sig := signature{short}
	// see also: parseSpec
	var cnt int
	export.WalkProperties(rtype, func(f *r.StructField, path []int) (done bool) {
		tags := tag.ReadTag(f.Tag)
		if _, ok := tags.Find("internal"); !ok {
			cnt++
			// write the selector:
			unlabeled := tags.Exists("unlabeled")
			if cnt == 1 && unlabeled {
				sig[0] += ":"
			}

			if cnt > 1 || !unlabeled {
				if f.Type.Kind() != r.Interface {
					//  write camel "fieldName:"
					name := firstRuneLower(f.Name)
					sig = sig.addSelector(name + ":")
				} else {
					slats := implementorsOf(f.Type)
					sig = sig.mulSelectors(slats)

				}
			}

		}
		return
	})
	//
	want := signature{"put:intoRec:atField:", "put:intoObj:atField:", "put:intoObjNamed:atField:"}
	if diff := pretty.Diff(sig, want); len(diff) > 0 {
		t.Fatal(sig)
	}
}

func typeName(name string, t r.Type) (ret string) {
	if len(name) > 0 {
		ret = name
	} else {
		ret = lang.Underscore(t.Name())
	}
	return
}

// return the command structs supported by the passed slot
func implementorsOf(slot r.Type) (ret []string) {
	for _, slats := range iffy.AllSlats {
		for _, slat := range slats {
			slat := r.TypeOf(slat)
			if slat.Implements(slot) {
				ret = append(ret, firstRuneLower(slat.Elem().Name()))
			}
		}
	}
	return
}

type signature []string

func (sig signature) addSelector(sel string) signature {
	for i, cnt := 0, len(sig); i < cnt; i++ {
		sig[i] = sig[i] + sel
	}
	return sig
}

func (sig signature) mulSelectors(sel []string) signature {
	var out signature
	for i, cnt := 0, len(sig); i < cnt; i++ {
		for _, sel := range sel {
			out = append(out, sig[i]+sel+":")
		}
	}
	return out
}

func shortName(t r.Type, name string) (ret string) {
	if len(name) > 0 {
		ret = name
	} else {
		ret = firstSyl(t.Name())
	}
	return
}

func firstRuneLower(s string) (ret string) {
	rs := []rune(s)
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}

func firstSyl(s string) (ret string) {
	var b strings.Builder
	for i, u := range s {
		if i == 0 {
			b.WriteRune(unicode.ToLower(u))
		} else if !unicode.IsUpper(u) {
			b.WriteRune(u)
		} else {
			break
		}
	}
	return b.String()
}
