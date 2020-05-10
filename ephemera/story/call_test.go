package story

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// test calling a pattern
// note: the pattern is undefined.
func TestDetermineNum(t *testing.T) {
	expect := core.DetermineNum{
		"factorial", &core.Parameters{[]*core.Parameter{{
			"num", &core.FromNum{
				&core.Number{3},
			}},
		}}}
	k, db := newTestDecoder(t)
	defer db.Close()
	if rule, e := imp_determine_num(k, _determine_num); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(rule, &expect); len(diff) != 0 {
		t.Fatal(diff)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select count() from eph_prog", 1)
		tables.WriteCsv(db, &buf, "select count() from eph_rule", 1)
		tables.WriteCsv(db, &buf, "select * from eph_pattern", 3)
		tables.WriteCsv(db, &buf, "select name, category from eph_named", 2)
		str := buf.String()
		if diff := pretty.Diff(str, lines(
			// eph_prog count
			// no programs b/c no container for the call into determine.
			"0",
			// eph_rule count
			// no rules b/c the pattern is called but not implemented.
			"0",
			// eph_pattern
			"1,2,3", // from NewPatternParam -> "determine num" takes a parameter that is from a number eval
			"1,1,4", // from NewPatternType -> "determine num" indicates factorial returns a number eval
			//
			"factorial,determine_num", // 1.
			"num,variable_name",       // 2.
			"number_eval,type",        // 3.
			"number_eval,type",        // 4.
		)); len(diff) > 0 {
			t.Fatal("mismatch", diff)
		} else {
			t.Log("ok", str)
		}
	}
}

// determine num of factorial where num = 3
var _determine_num = map[string]interface{}{
	"type": "determine_num",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "factorial",
		},
		"$PARAMETERS": map[string]interface{}{
			"type": "parameters",
			"value": map[string]interface{}{
				"$PARAMS": []interface{}{
					map[string]interface{}{
						"type": "parameter",
						"value": map[string]interface{}{
							"$FROM": map[string]interface{}{
								"type": "assignment",
								"value": map[string]interface{}{
									"type": "assign_num",
									"value": map[string]interface{}{
										"$VAL": map[string]interface{}{
											"type": "number_eval",
											"value": map[string]interface{}{
												"type": "num_value",
												"value": map[string]interface{}{
													"$NUM": map[string]interface{}{
														"type":  "number",
														"value": 3.0, // json numbers are float64
													}}}}}},
							},
							"$NAME": map[string]interface{}{
								"type":  "variable_name",
								"value": "num",
							}}}}}}},
}
