package story

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/kr/pretty"
)

// test importing using structures
// ( so that we can log failures in unread or missing elements )
func TestRefresh(t *testing.T) {
	const filePath = "/Users/ionous/Dev/go/src/github.com/ionous/iffy/stories/shared/plurals.if"
	if p, e := readJson(filePath); e != nil {
		t.Fatal(e)
	} else {
		var ds reader.Dilemmas
		dec := decode.NewDecoderReporter("plurals.if", ds.Report)
		dec.AddDefaultCallbacks(Slats)
		if i, e := dec.ReadSpec(p); e != nil {
			t.Fatal(e)
		} else {
			reader.PrintDilemmas(log.Writer(), ds)
			t.Log(pretty.Sprint(i))
		}
	}
}

func readJson(filePath string) (ret reader.Map, err error) {
	if f, e := os.Open(filePath); e != nil {
		err = e
	} else {
		defer f.Close()
		dec := json.NewDecoder(f)
		if e := dec.Decode(&ret); e != nil && e != io.EOF {
			err = e
		}
	}
	return
}

var Slots = []composer.Slot{
	{Type: (*StoryStatement)(nil)},
}

var Slats = []composer.Slat{
	(*Story)(nil),
	(*Paragraph)(nil),
	// story statements:
	(*KindsPossessProperties)(nil),
}

type Story struct {
	Paragraph []*Paragraph
}
type Paragraph struct {
	StoryStatement []StoryStatement
}
type StoryStatement interface {
	StoryStatement() // marker interface
}
type KindsPossessProperties struct {
	PluralKinds
	Determiner
	PropertyPhrase struct {
		*PrimitivePhrase
		*AspectPhrase
	}
	// PropertyPhrase: {primitive_phrase} or {aspect_phrase}
}

type PrimitivePhrase struct {
	PrimitiveType
	Property
}
type AspectPhrase struct {
	Aspect
	*OptionalProperty
}

type OptionalProperty struct {
	Property
}

type PrimitiveType int

//go:generate stringer -type=PrimitiveType
const (
	PrimitiveTypeNumber PrimitiveType = iota
	PrimitiveTypeText
	PrimitiveTypeBool
)

type Aspect string
type PluralKinds string
type Determiner string
type Property string

func (*Story) Compose() (_ composer.Spec)                  { return }
func (*Paragraph) Compose() (_ composer.Spec)              { return }
func (*KindsPossessProperties) Compose() (_ composer.Spec) { return }
func (*KindsPossessProperties) StoryStatement()            {}
