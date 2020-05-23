package story

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
)

func TestImportSequence(t *testing.T) {
	k, db := newTestDecoder(t)
	defer db.Close()
	if cmd, e := imp_cycle_text(k, _cycle_text); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(cmd, &core.CycleText{
		Sequence: core.Sequence{
			Seq: "autoimp1",
			Parts: []rt.TextEval{
				&core.Text{Text: "a"},
				&core.Text{Text: "b"},
				&core.Text{Text: "c"},
			},
		},
	}); len(diff) > 0 {
		t.Fatal(pretty.Print(cmd))
	}
}

var _cycle_text = map[string]interface{}{
	"type": "cycle_text",
	"value": map[string]interface{}{
		"$PARTS": []interface{}{
			map[string]interface{}{
				"type": "text_eval",
				"value": map[string]interface{}{
					"type": "text_value",
					"value": map[string]interface{}{
						"$TEXT": map[string]interface{}{
							"type":  "lines",
							"value": "a",
						}}}},
			map[string]interface{}{
				"type": "text_eval",
				"value": map[string]interface{}{
					"type": "text_value",
					"value": map[string]interface{}{
						"$TEXT": map[string]interface{}{
							"type":  "lines",
							"value": "b",
						}}}},
			map[string]interface{}{
				"type": "text_eval",
				"value": map[string]interface{}{
					"type": "text_value",
					"value": map[string]interface{}{
						"$TEXT": map[string]interface{}{
							"type":  "lines",
							"value": "c",
						}}}}}},
}
