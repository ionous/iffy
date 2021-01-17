package debug

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/ephemera/story"
	"github.com/ionous/iffy/rt"
)

var FactorialStory = &story.Story{
	Paragraph: &[]story.Paragraph{{
		StoryStatement: &[]story.StoryStatement{
			&story.TestStatement{
				TestName: story.TestName{
					Str: "factorial",
				},
				Test: &story.TestOutput{
					Lines: story.Lines{
						Str: "6",
					}}},
			&story.TestRule{
				TestName: story.TestName{
					Str: "factorial",
				},
				Hook: story.ProgramHook{
					Opt: &story.Activity{
						Exe: []rt.Execute{
							&core.Say{
								Text: &core.PrintNum{
									Num: &pattern.Determine{
										Pattern: "factorial",
										Arguments: &core.Arguments{
											Args: []*core.Argument{
												&core.Argument{
													Name: "num",
													From: &core.FromNum{
														Val: &core.Number{Num: 3}},
												}},
										}},
								}},
						}},
				}},
			&story.PatternDecl{
				Name: story.PatternName{
					Str: "factorial",
				},
				Type: story.PatternType{
					Opt: &story.PatternedActivity{
						Str: "$ACTIVITY"}},
				Optvars: &story.PatternVariablesTail{
					VariableDecl: []story.VariableDecl{numberDecl}},
			},
			&story.PatternActions{
				Name: story.PatternName{
					Str: "factorial",
				},
				PatternReturn: &story.PatternReturn{Result: numberDecl},
				PatternRules: story.PatternRules{
					PatternRule: &[]story.PatternRule{{
						Guard: &core.Always{},
						Hook: story.ProgramHook{
							Opt: &story.Activity{[]rt.Execute{
								&core.Assign{
									Var: core.Variable{Str: "num"},
									From: &core.FromNum{&core.ProductOf{
										A: &core.Var{Name: "num"},
										B: &core.DiffOf{
											A: &core.Var{Name: "num"},
											B: &core.Number{Num: 1}},
									}},
								},
							}},
						}},
					}}},
			&story.PatternActions{
				Name: story.PatternName{
					Str: "factorial",
				},
				PatternReturn: &story.PatternReturn{Result: numberDecl},
				PatternRules: story.PatternRules{
					PatternRule: &[]story.PatternRule{{
						Guard: &core.CompareNum{
							A:  &core.Var{Name: "num"},
							Is: &core.EqualTo{},
							B:  &core.Number{}},
						Hook: story.ProgramHook{
							Opt: &story.Activity{[]rt.Execute{
								&core.Assign{
									Var:  core.Variable{Str: "num"},
									From: &core.FromNum{&core.Number{Num: 1}},
								}},
							}}},
					}},
			}},
	}},
}

var numberDecl = story.VariableDecl{
	An: story.Determiner{
		Str: "a",
	},
	Name: story.VariableName{
		Variable: core.Variable{
			Str: "num",
		}},
	Type: story.VariableType{
		Opt: &story.PrimitiveType{
			Str: "$NUMBER",
		}},
}
