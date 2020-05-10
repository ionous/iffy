package debug

import (
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/rt"
)

var SayTest = check.TestOutput{
	"hello", []rt.Execute{
		&core.Choose{
			If: &core.Bool{true},
			True: []rt.Execute{&core.Say{
				Text: &core.Text{"hello"},
			}},
			False: []rt.Execute{&core.Say{
				Text: &core.Text{"goodbye"},
			}},
		}},
}

var SayStory = reader.Map{
	"type": "test_output",
	"value": map[string]interface{}{
		"$LINES": map[string]interface{}{
			"type":  "lines",
			"value": "hello",
		},
		"$GO": []interface{}{
			map[string]interface{}{
				"type": "execute",
				"value": map[string]interface{}{
					"type": "choose",
					"value": map[string]interface{}{
						"$FALSE": []interface{}{
							map[string]interface{}{
								"type": "execute",
								"value": map[string]interface{}{
									"type": "say_text",
									"value": map[string]interface{}{
										"$TEXT": map[string]interface{}{
											"type": "text_eval",
											"value": map[string]interface{}{
												"type": "text_value",
												"value": map[string]interface{}{
													"$TEXT": map[string]interface{}{
														"type":  "lines",
														"value": "goodbye",
													}}}}}}}},
						"$IF": map[string]interface{}{
							"type": "bool_eval",
							"value": map[string]interface{}{
								"type": "bool_value",
								"value": map[string]interface{}{
									"$BOOL": map[string]interface{}{
										"type":  "bool",
										"value": "$TRUE",
									}}}},
						"$TRUE": []interface{}{
							map[string]interface{}{
								"type": "execute",
								"value": map[string]interface{}{
									"type": "say_text",
									"value": map[string]interface{}{
										"$TEXT": map[string]interface{}{
											"type": "text_eval",
											"value": map[string]interface{}{
												"type": "text_value",
												"value": map[string]interface{}{
													"$TEXT": map[string]interface{}{
														"type":  "lines",
														"value": "hello",
													}}}}}}}}}}}}},
}
