package debug

var FactorialStory = map[string]interface{}{
	"type": "story",
	"value": map[string]interface{}{
		"$PARAGRAPH": []interface{}{
			map[string]interface{}{
				"type": "paragraph",
				"value": map[string]interface{}{
					"$STORY_STATEMENT": []interface{}{
						map[string]interface{}{
							"type":  "story_statement",
							"value": FactorialTestOutput,
						},
						map[string]interface{}{
							"type":  "story_statement",
							"value": FactorialTestRule,
						},
						map[string]interface{}{
							"type":  "story_statement",
							"value": FactorialPatternDecl,
						},
						map[string]interface{}{
							"type":  "story_statement",
							"value": FactorialSubtract,
						},
						map[string]interface{}{
							"type":  "story_statement",
							"value": FactorialZero,
						}}}}}},
}

var FactoialNamedTest = map[string]interface{}{
	"type":  "test_name",
	"value": "factorial",
}

var FactorialTestOutput = map[string]interface{}{
	"type": "test_statement",
	"value": map[string]interface{}{
		"$TEST_NAME": FactoialNamedTest,
		"$TEST": map[string]interface{}{
			"type": "testing",
			"value": map[string]interface{}{
				"type": "test_output",
				"value": map[string]interface{}{
					"$LINES": map[string]interface{}{
						"type":  "lines",
						"value": "6",
					}}}}},
}

var FactorialTestRule = map[string]interface{}{
	"type": "test_rule",
	"value": map[string]interface{}{
		"$TEST_NAME": FactoialNamedTest,
		"$HOOK": map[string]interface{}{
			"type": "program_hook",
			"value": map[string]interface{}{
				"$ACTIVITY": map[string]interface{}{
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
												"type": "print_num",
												"value": map[string]interface{}{
													"$NUM": map[string]interface{}{
														"type":  "number_eval",
														"value": FactorialDetermineNum,
													}}}}}}}}}}}}},
}

// determine num of factorial where num = 3
var FactorialDetermineNum = map[string]interface{}{
	"type": "determine_num",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "factorial",
		},
		"$ARGUMENTS": map[string]interface{}{
			"type": "arguments",
			"value": map[string]interface{}{
				"$ARGS": []interface{}{
					map[string]interface{}{
						"type": "argument",
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
														"value": 3.0,
													}}}}}}},
							"$NAME": map[string]interface{}{
								"type":  "variable_name",
								"value": "num",
							}}}}}}},
}

var FactorialPatternDecl = map[string]interface{}{
	"type": "pattern_decl",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "factorial",
		},
		"$TYPE": map[string]interface{}{
			"type": "pattern_type",
			"value": map[string]interface{}{
				"$VALUE": map[string]interface{}{
					"type": "variable_type",
					"value": map[string]interface{}{
						"$PRIMITIVE": map[string]interface{}{
							"type":  "primitive_type",
							"value": "$NUMBER",
						}}}}},
		"$OPTVARS": map[string]interface{}{
			"type": "pattern_variables_tail",
			"value": map[string]interface{}{
				"$VARIABLE_DECL": []interface{}{
					map[string]interface{}{
						"type": "variable_decl",
						"value": map[string]interface{}{
							"$TYPE": map[string]interface{}{
								"type": "variable_type",
								"value": map[string]interface{}{
									"$PRIMITIVE": map[string]interface{}{
										"type":  "primitive_type",
										"value": "$NUMBER",
									}}},
							"$NAME": map[string]interface{}{
								"type":  "variable_name",
								"value": "num",
							}}}}}}},
}

var FactorialZero = map[string]interface{}{
	"type": "pattern_actions",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "factorial",
		},
		"$PATTERN_RULES": map[string]interface{}{
			"type": "pattern_rules",
			"value": map[string]interface{}{
				"$PATTERN_RULE": []interface{}{
					map[string]interface{}{
						"type": "pattern_rule",
						"value": map[string]interface{}{
							"$HOOK": map[string]interface{}{
								"type": "program_hook",
								"value": map[string]interface{}{
									"$RESULT": map[string]interface{}{
										"type": "program_return",
										"value": map[string]interface{}{
											"$RESULT": map[string]interface{}{
												"type": "program_result",
												"value": map[string]interface{}{
													"$PRIMITIVE": map[string]interface{}{
														"type": "primitive_func",
														"value": map[string]interface{}{
															"$NUMBER_EVAL": map[string]interface{}{
																"type": "number_eval",
																"value": map[string]interface{}{
																	"type": "num_value",
																	"value": map[string]interface{}{
																		"$NUM": map[string]interface{}{
																			"type":  "number",
																			"value": 1.0,
																		}}}}}}}}}}}},
							"$GUARD": map[string]interface{}{
								"type": "bool_eval",
								"value": map[string]interface{}{
									"type": "compare_num",
									"value": map[string]interface{}{
										"$A": map[string]interface{}{
											"type": "number_eval",
											"value": map[string]interface{}{
												"type": "get_var",
												"value": map[string]interface{}{
													"$NAME": FactorialNum,
												}}},
										"$IS": map[string]interface{}{
											"type": "comparator",
											"value": map[string]interface{}{
												"type":  "equal",
												"value": map[string]interface{}{},
											}},
										"$B": map[string]interface{}{
											"type": "number_eval",
											"value": map[string]interface{}{
												"type": "num_value",
												"value": map[string]interface{}{
													"$NUM": map[string]interface{}{
														"type":  "number",
														"value": 0.0,
													}}}}}}}}}}}}},
}

var FactorialNum = map[string]interface{}{
	"type": "text_eval",
	"value": map[string]interface{}{
		"type": "text_value",
		"value": map[string]interface{}{
			"$TEXT": map[string]interface{}{
				"type":  "text",
				"value": "num",
			},
		},
	},
}

var FactorialSubtract = map[string]interface{}{
	"type": "pattern_actions",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "factorial",
		},
		"$PATTERN_RULES": map[string]interface{}{
			"type": "pattern_rules",
			"value": map[string]interface{}{
				"$PATTERN_RULE": []interface{}{
					map[string]interface{}{
						"type": "pattern_rule",
						"value": map[string]interface{}{
							"$GUARD": map[string]interface{}{
								"type": "bool_eval",
								"value": map[string]interface{}{
									"type":  "always",
									"value": map[string]interface{}{},
								}},
							"$HOOK": map[string]interface{}{
								"type": "program_hook",
								"value": map[string]interface{}{
									"$RESULT": map[string]interface{}{
										"type": "program_return",
										"value": map[string]interface{}{
											"$RESULT": map[string]interface{}{
												"type": "program_result",
												"value": map[string]interface{}{
													"$PRIMITIVE": map[string]interface{}{
														"type": "primitive_func",
														"value": map[string]interface{}{
															"$NUMBER_EVAL": map[string]interface{}{
																"type": "number_eval",
																"value": map[string]interface{}{
																	"type": "product_of",
																	"value": map[string]interface{}{
																		"$A": map[string]interface{}{
																			"type": "number_eval",
																			"value": map[string]interface{}{
																				"type": "get_var",
																				"value": map[string]interface{}{
																					"$NAME": FactorialNum,
																				}}},
																		"$B": map[string]interface{}{
																			"type": "number_eval",
																			"value": map[string]interface{}{
																				"type": "diff_of",
																				"value": map[string]interface{}{
																					"$A": map[string]interface{}{
																						"type": "number_eval",
																						"value": map[string]interface{}{
																							"type": "get_var",
																							"value": map[string]interface{}{
																								"$NAME": FactorialNum,
																							}}},
																					"$B": map[string]interface{}{
																						"type": "number_eval",
																						"value": map[string]interface{}{
																							"type": "num_value",
																							"value": map[string]interface{}{
																								"$NUM": map[string]interface{}{
																									"type":  "number",
																									"value": 1.0,
																								}}}}}}}}}}}}}}}}}}}}}}}},
}
