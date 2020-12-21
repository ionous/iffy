package story

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/test/testdb"
	"github.com/kr/pretty"
)

func TestImportSequence(t *testing.T) {
	k, db := newImporter(t, testdb.Memory)
	defer db.Close()
	if cmd, e := k.decoder.ReadSpec(_cycle_text); e != nil {
		t.Fatal("failed to read sequence", e)
	} else if diff := pretty.Diff(cmd, &core.CycleText{
		Sequence: core.Sequence{
			Seq: "seq_1", // COUNTER:#
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
							"type":  "text",
							"value": "a",
						}}}},
			map[string]interface{}{
				"type": "text_eval",
				"value": map[string]interface{}{
					"type": "text_value",
					"value": map[string]interface{}{
						"$TEXT": map[string]interface{}{
							"type":  "text",
							"value": "b",
						}}}},
			map[string]interface{}{
				"type": "text_eval",
				"value": map[string]interface{}{
					"type": "text_value",
					"value": map[string]interface{}{
						"$TEXT": map[string]interface{}{
							"type":  "text",
							"value": "c",
						}}}}}},
}
