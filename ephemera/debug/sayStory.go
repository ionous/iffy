package debug

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/reader"
)

var SayTest = core.NewActivity(
	&core.Choose{
		If: &core.Bool{true},
		True: core.NewActivity(&core.Say{
			Text: &core.Text{"hello"},
		}),
		False: core.NewActivity(&core.Say{
			Text: &core.Text{"goodbye"},
		}),
	})

var SayStory = reader.Map{
	"type": "activity",
	"value": map[string]interface{}{
		"$EXE": []interface{}{
			map[string]interface{}{
				"type": "execute",
				"value": map[string]interface{}{
					"type": "choose",
					"value": map[string]interface{}{
						"$FALSE": map[string]interface{}{
							"type": "activity",
							"value": map[string]interface{}{
								"$EXE": []interface{}{
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
																"type":  "text",
																"value": "goodbye",
															}}}}}}}}}},
						"$IF": map[string]interface{}{
							"type": "bool_eval",
							"value": map[string]interface{}{
								"type": "bool_value",
								"value": map[string]interface{}{
									"$BOOL": map[string]interface{}{
										"type":  "bool",
										"value": "$TRUE",
									}}}},
						"$TRUE": map[string]interface{}{
							"type": "activity",
							"value": map[string]interface{}{
								"$EXE": []interface{}{
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
																"type":  "text",
																"value": "hello",
															}}}}}}}}}}}}}}},
}
