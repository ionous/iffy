// Code generated by "makeops"; edit at your own risk.
package story

import (
  "github.com/ionous/iffy/dl/composer"
  "github.com/ionous/iffy/ephemera/reader"
  "github.com/ionous/iffy/rt"
)

// An requires a user-specified string.
type An struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *An) String() string {
  return op.Str
}

func (*An) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    "$A": "a",
  }
}

func (*An) Compose() composer.Spec {
  return composer.Spec{
    Name: "an",
    Spec: "{a} or {an}",
  }
}

// AreAn requires a user-specified string.
type AreAn struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *AreAn) String() string {
  return op.Str
}

func (*AreAn) Choices() (closed bool, choices map[string]string) {
  return true, map[string]string{
    "$ARE": "are", "$AREA": "area", "$AREAN": "arean", "$IS": "is", "$ISA": "isa", "$ISAN": "isan",
  }
}

func (*AreAn) Compose() composer.Spec {
  return composer.Spec{
    Name: "are_an",
    Spec: "{are}, {are a%area}, {are an%arean}, {is}, {is a%isa}, {is an%isan}",
  }
}

// AreBeing requires a user-specified string.
type AreBeing struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *AreBeing) String() string {
  return op.Str
}

func (*AreBeing) Choices() (closed bool, choices map[string]string) {
  return true, map[string]string{
    "$ARE": "are", "$IS": "is",
  }
}

func (*AreBeing) Compose() composer.Spec {
  return composer.Spec{
    Name: "are_being",
    Spec: "{are} or {is}",
  }
}

// AreEither requires a user-specified string.
type AreEither struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *AreEither) String() string {
  return op.Str
}

func (*AreEither) Choices() (closed bool, choices map[string]string) {
  return true, map[string]string{
    "$CANBE": "canbe", "$EITHER": "either",
  }
}

func (*AreEither) Compose() composer.Spec {
  return composer.Spec{
    Name: "are_either",
    Spec: "{can be%canbe} {are either%either}",
  }
}

// Argument requires various parameters.
type Argument struct {
  At reader.Position `if:"internal"`
  Name VariableName
  From Assignment
}

func (*Argument) Compose() composer.Spec {
  return composer.Spec{
    Name: "argument",
    Group: "patterns",
  }
}

// Arguments requires various parameters.
type Arguments struct {
  At reader.Position `if:"internal"`
  Args []Argument
}

func (*Arguments) Compose() composer.Spec {
  return composer.Spec{
    Name: "arguments",
    Group: "patterns",
  }
}

// Aspect requires a user-specified string.
type Aspect struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Aspect) String() string {
  return op.Str
}

func (*Aspect) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*Aspect) Compose() composer.Spec {
  return composer.Spec{
    Name: "aspect",
    Spec: "{aspect}",
  }
}

// AspectPhrase requires various parameters.
type AspectPhrase struct {
  At reader.Position `if:"internal"`
  Aspect Aspect
  OptionalProperty *OptionalProperty
}

func (*AspectPhrase) Compose() composer.Spec {
  return composer.Spec{
    Name: "aspect_phrase",
    Spec: "{aspect} {?optional_property}",
  }
}

// AspectTraits requires various parameters.
type AspectTraits struct {
  At reader.Position `if:"internal"`
  Aspect Aspect
  TraitPhrase TraitPhrase
}

func (*AspectTraits) Compose() composer.Spec {
  return composer.Spec{
    Name: "aspect_traits",
    Spec: "{aspect} {trait_phrase}",
  }
}

// Bool requires a user-specified string.
type Bool struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Bool) String() string {
  return op.Str
}

func (*Bool) Choices() (closed bool, choices map[string]string) {
  return true, map[string]string{
    "$TRUE": "true", "$FALSE": "false",
  }
}

func (*Bool) Compose() composer.Spec {
  return composer.Spec{
    Name: "bool",
    Spec: "{true} or {false}",
  }
}

// BoxedNumber requires various parameters.
type BoxedNumber struct {
  At reader.Position `if:"internal"`
  Number Number
}

func (*BoxedNumber) Compose() composer.Spec {
  return composer.Spec{
    Name: "boxed_number",
    Spec: "{number}",
  }
}

// BoxedText requires various parameters.
type BoxedText struct {
  At reader.Position `if:"internal"`
  Text Text
}

func (*BoxedText) Compose() composer.Spec {
  return composer.Spec{
    Name: "boxed_text",
    Spec: "{text}",
  }
}

// Certainties requires various parameters.
type Certainties struct {
  At reader.Position `if:"internal"`
  PluralKinds PluralKinds
  AreBeing AreBeing
  Certainty Certainty
  Trait Trait
}

func (*Certainties) Compose() composer.Spec {
  return composer.Spec{
    Name: "certainties",
    Spec: "{plural_kinds} {are_being} {certainty} {trait}.",
  }
}

// Certainty requires a user-specified string.
type Certainty struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Certainty) String() string {
  return op.Str
}

func (*Certainty) Choices() (closed bool, choices map[string]string) {
  return true, map[string]string{
    "$USUALLY": "usually", "$ALWAYS": "always", "$SELDOM": "seldom", "$NEVER": "never",
  }
}

func (*Certainty) Compose() composer.Spec {
  return composer.Spec{
    Name: "certainty",
    Desc: `Certainty: Whether an trait applies to a kind of noun.`,
    Spec: "{usually}, {always}, {seldom}, or {never}",
  }
}

// Comment requires various parameters.
type Comment struct {
  At reader.Position `if:"internal"`
  Lines Lines
}

func (*Comment) Compose() composer.Spec {
  return composer.Spec{
    Name: "comment",
    Desc: `Comment: Information about the story not used by the story.`,
    Spec: "Note: {comment%lines}",
  }
}

// Comments requires various parameters.
type Comments struct {
  At reader.Position `if:"internal"`
  Lines Lines
}

func (*Comments) Compose() composer.Spec {
  return composer.Spec{
    Name: "comments",
    Spec: "{lines|quote}",
  }
}

// CycleText requires various parameters.
type CycleText struct {
  At reader.Position `if:"internal"`
  Parts []rt.TextEval
}

func (*CycleText) Compose() composer.Spec {
  return composer.Spec{
    Name: "cycle_text",
    Desc: `Cycle text: When called multiple times, returns each of its inputs in turn.`,
    Group: "cycle",
  }
}

// DetermineAct requires various parameters.
type DetermineAct struct {
  At reader.Position `if:"internal"`
  Name PatternName
  Arguments *Arguments
}

func (*DetermineAct) Compose() composer.Spec {
  return composer.Spec{
    Name: "determine_act",
    Desc: `Determine an activity`,
    Group: "patterns",
  }
}

// DetermineBool requires various parameters.
type DetermineBool struct {
  At reader.Position `if:"internal"`
  Name PatternName
  Arguments *Arguments
}

func (*DetermineBool) Compose() composer.Spec {
  return composer.Spec{
    Name: "determine_bool",
    Desc: `Determine a true/false value`,
    Group: "patterns",
  }
}

// DetermineNum requires various parameters.
type DetermineNum struct {
  At reader.Position `if:"internal"`
  Name PatternName
  Arguments *Arguments
}

func (*DetermineNum) Compose() composer.Spec {
  return composer.Spec{
    Name: "determine_num",
    Desc: `Determine a number`,
    Group: "patterns",
  }
}

// DetermineNumList requires various parameters.
type DetermineNumList struct {
  At reader.Position `if:"internal"`
  Name PatternName
  Arguments *Arguments
}

func (*DetermineNumList) Compose() composer.Spec {
  return composer.Spec{
    Name: "determine_num_list",
    Desc: `Determine a list of numbers`,
    Group: "patterns",
  }
}

// DetermineText requires various parameters.
type DetermineText struct {
  At reader.Position `if:"internal"`
  Name PatternName
  Arguments *Arguments
}

func (*DetermineText) Compose() composer.Spec {
  return composer.Spec{
    Name: "determine_text",
    Desc: `Determine some text`,
    Group: "patterns",
  }
}

// DetermineTextList requires various parameters.
type DetermineTextList struct {
  At reader.Position `if:"internal"`
  Name PatternName
  Arguments *Arguments
}

func (*DetermineTextList) Compose() composer.Spec {
  return composer.Spec{
    Name: "determine_text_list",
    Desc: `Determine a list of text`,
    Group: "patterns",
  }
}

// Determiner requires a user-specified string.
type Determiner struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Determiner) String() string {
  return op.Str
}

func (*Determiner) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    "$A": "a", "$AN": "an", "$THE": "the", "$OUR": "our",
  }
}

func (*Determiner) Compose() composer.Spec {
  return composer.Spec{
    Name: "determiner",
    Desc: `Determiners: modify a word they are associated to designate specificity or, sometimes, a count. For instance: "some" fish hooks, "a" pineapple, "75" triangles, "our" Trevor.`,
    Spec: "{a}, {an}, {the}, {our}, or {other determiner%determiner}",
  }
}

// ExtType swaps between various options
type ExtType struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*ExtType) Compose() composer.Spec {
  return composer.Spec{
    Name: "ext_type",
    Spec: "a list of {numbers:number_list}, a list of {text%text_list}, a {record:record_type} or a list of {records:record_list}.",
  }
}

func (*ExtType) Choices() map[string]interface{} {
  return map[string]interface{} {
    "numbers": (*NumberList)(nil),
    "text_list": (*TextList)(nil),
    "record": (*RecordType)(nil),
    "records": (*RecordList)(nil),
  }
}

// KindOfNoun requires various parameters.
type KindOfNoun struct {
  At reader.Position `if:"internal"`
  AreAn AreAn
  Trait *[]Trait
  Kind SingularKind
  NounRelation *NounRelation
}

func (*KindOfNoun) Compose() composer.Spec {
  return composer.Spec{
    Name: "kind_of_noun",
    Spec: "{are_an} {trait*trait|comma-and} kind of {kind:singular_kind} {?noun_relation}",
  }
}

// KindsOfAspect requires various parameters.
type KindsOfAspect struct {
  At reader.Position `if:"internal"`
  Aspect Aspect
}

func (*KindsOfAspect) Compose() composer.Spec {
  return composer.Spec{
    Name: "kinds_of_aspect",
    Spec: "{aspect} is a kind of value.",
  }
}

// KindsOfKind requires various parameters.
type KindsOfKind struct {
  At reader.Position `if:"internal"`
  PluralKinds PluralKinds
  SingularKind SingularKind
}

func (*KindsOfKind) Compose() composer.Spec {
  return composer.Spec{
    Name: "kinds_of_kind",
    Spec: "{plural_kinds} are a kind of {singular_kind}.",
  }
}

// KindsPossessProperties requires various parameters.
type KindsPossessProperties struct {
  At reader.Position `if:"internal"`
  PluralKinds PluralKinds
  Determiner Determiner
  PropertyPhrase PropertyPhrase
}

func (*KindsPossessProperties) Compose() composer.Spec {
  return composer.Spec{
    Name: "kinds_possess_properties",
    Spec: "{plural_kinds} have {determiner} {property_phrase}.",
  }
}

// Lede requires various parameters.
type Lede struct {
  At reader.Position `if:"internal"`
  Nouns []NamedNoun
  NounPhrase NounPhrase
}

func (*Lede) Compose() composer.Spec {
  return composer.Spec{
    Name: "lede",
    Desc: `Leading statement: Describes one or more nouns.`,
    Spec: "{nouns+named_noun|comma-and} {noun_phrase}.",
  }
}

// Lines requires a user-specified string.
type Lines struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Lines) String() string {
  return op.Str
}

func (*Lines) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*Lines) Compose() composer.Spec {
  return composer.Spec{
    Name: "lines",
    Desc: `Lines: A sequence of characters of any length spanning multiple lines. Paragraphs are a prime example. Generally lines are some piece of the story that will be displayed to the player.
See also: text.`,
    Spec: "{lines}",
  }
}

// LocalDecl requires various parameters.
type LocalDecl struct {
  At reader.Position `if:"internal"`
  VariableDecl VariableDecl
  ProgramResult ProgramResult
}

func (*LocalDecl) Compose() composer.Spec {
  return composer.Spec{
    Name: "local_decl",
    Desc: `Local: local variables can use the parameters of a pattern to compute temporary values.`,
    Spec: "Where {variable_decl} = {value%program_result}",
  }
}

// NamedNoun requires various parameters.
type NamedNoun struct {
  At reader.Position `if:"internal"`
  Determiner Determiner
  Name NounName
}

func (*NamedNoun) Compose() composer.Spec {
  return composer.Spec{
    Name: "named_noun",
    Spec: "{determiner} {name:noun_name}",
  }
}

// NounAssignment requires various parameters.
type NounAssignment struct {
  At reader.Position `if:"internal"`
  Property Property
  Nouns []NamedNoun
  Lines Lines
}

func (*NounAssignment) Compose() composer.Spec {
  return composer.Spec{
    Name: "noun_assignment",
    Desc: `Noun assignment: Assign text. Gives a noun one or more lines of text.`,
    Spec: "The {property} of {nouns+named_noun} is {the text%lines|summary}",
  }
}

// NounName requires a user-specified string.
type NounName struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *NounName) String() string {
  return op.Str
}

func (*NounName) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*NounName) Compose() composer.Spec {
  return composer.Spec{
    Name: "noun_name",
    Desc: `Noun name: Some specific person, place, or thing; or, more rarely, a kind. Proper names are usually capitalized:  For example, maybe: 'Haruki', 'Jane', or 'Toronto'.
        Common names are usually not capitalized. For example, maybe: 'table', 'chair', or 'dog park'.
        A set of duplicate object uses their kind. For instance: twelve 'cats'.`,
    Spec: "{noun_name}",
  }
}

// NounPhrase swaps between various options
type NounPhrase struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*NounPhrase) Compose() composer.Spec {
  return composer.Spec{
    Name: "noun_phrase",
    Spec: "{kind_of_noun}, {noun_traits}, or {noun_relation}",
  }
}

func (*NounPhrase) Choices() map[string]interface{} {
  return map[string]interface{} {
    "kind_of_noun": (*KindOfNoun)(nil),
    "noun_traits": (*NounTraits)(nil),
    "noun_relation": (*NounRelation)(nil),
  }
}

// NounRelation requires various parameters.
type NounRelation struct {
  At reader.Position `if:"internal"`
  AreBeing *AreBeing
  Relation Relation
  Nouns []NamedNoun
}

func (*NounRelation) Compose() composer.Spec {
  return composer.Spec{
    Name: "noun_relation",
    Spec: "{?are_being} {relation} {nouns+named_noun|comma-and}",
  }
}

// NounStatement requires various parameters.
type NounStatement struct {
  At reader.Position `if:"internal"`
  Lede Lede
  Tail *[]Tail
  Summary *Summary
}

func (*NounStatement) Compose() composer.Spec {
  return composer.Spec{
    Name: "noun_statement",
    Desc: `Noun statement: Describes people, places, or things.`,
    Spec: "{:lede} {*tail} {?summary}",
  }
}

// NounTraits requires various parameters.
type NounTraits struct {
  At reader.Position `if:"internal"`
  AreBeing AreBeing
  Trait []Trait
}

func (*NounTraits) Compose() composer.Spec {
  return composer.Spec{
    Name: "noun_traits",
    Spec: "{are_being} {trait+trait|comma-and}",
  }
}

// NounType requires various parameters.
type NounType struct {
  At reader.Position `if:"internal"`
  An An
  Kinds PluralKinds
}

func (*NounType) Compose() composer.Spec {
  return composer.Spec{
    Name: "noun_type",
    Spec: "{an} {kind of%kinds:plural_kinds} noun",
  }
}

// Number requires a user-specified number.
type Number struct {
  At reader.Position `if:"internal"`
  Val float64
}


func (*Number) Num() (closed bool, choices []float64) {
    return false, []float64{
    
  }
}

func (*Number) Compose() composer.Spec {
  return composer.Spec{
    Name: "number",
  }
}

// NumberList requires a user-specified string.
type NumberList struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *NumberList) String() string {
  return op.Str
}

func (*NumberList) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*NumberList) Compose() composer.Spec {
  return composer.Spec{
    Name: "number_list",
    Spec: "{a list of numbers%number_list}",
  }
}

// ObjectFunc requires various parameters.
type ObjectFunc struct {
  At reader.Position `if:"internal"`
  Name rt.TextEval
}

func (*ObjectFunc) Compose() composer.Spec {
  return composer.Spec{
    Name: "object_func",
    Spec: "an object named {name:text_eval}",
  }
}

// ObjectType requires various parameters.
type ObjectType struct {
  At reader.Position `if:"internal"`
  An An
  Kind SingularKind
}

func (*ObjectType) Compose() composer.Spec {
  return composer.Spec{
    Name: "object_type",
    Spec: "{an} {kind of%kind:singular_kind}",
  }
}

// OptionalProperty requires various parameters.
type OptionalProperty struct {
  At reader.Position `if:"internal"`
  Property Property
}

func (*OptionalProperty) Compose() composer.Spec {
  return composer.Spec{
    Name: "optional_property",
    Spec: "called {property}",
  }
}

// Paragraph requires various parameters.
type Paragraph struct {
  At reader.Position `if:"internal"`
  StoryStatement *[]StoryStatement
}

func (*Paragraph) Compose() composer.Spec {
  return composer.Spec{
    Name: "paragraph",
    Desc: `Phrases`,
    Spec: "{*story_statement}",
  }
}

// PatternActions requires various parameters.
type PatternActions struct {
  At reader.Position `if:"internal"`
  Name PatternName
  PatternLocals *PatternLocals
  PatternRules PatternRules
}

func (*PatternActions) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_actions",
    Desc: `Pattern actions: Actions to take when using a pattern.`,
    Spec: "To {pattern name%name:pattern_name}: {?pattern_locals} {pattern_rules}",
  }
}

// PatternDecl requires various parameters.
type PatternDecl struct {
  At reader.Position `if:"internal"`
  Name PatternName
  Type PatternType
  Optvars *PatternVariablesTail
  About *Comments
}

func (*PatternDecl) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_decl",
    Desc: `Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it. Each function in a given pattern has "guards" which determine whether the function applies in a particular situtation.`,
    Spec: "The pattern {name:pattern_name|quote} determines {type:pattern_type}. {optvars?pattern_variables_tail} {about?comments}",
  }
}

// PatternLocals requires various parameters.
type PatternLocals struct {
  At reader.Position `if:"internal"`
  LocalDecl *[]LocalDecl
}

func (*PatternLocals) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_locals",
    Spec: "{*local_decl}",
  }
}

// PatternName requires a user-specified string.
type PatternName struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *PatternName) String() string {
  return op.Str
}

func (*PatternName) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*PatternName) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_name",
    Spec: "{pattern_name}",
  }
}

// PatternRule requires various parameters.
type PatternRule struct {
  At reader.Position `if:"internal"`
  Guard rt.BoolEval
  Hook ProgramHook
}

func (*PatternRule) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_rule",
    Desc: `Rule`,
    Spec: "When {conditions are met%guard:bool_eval}, then: {do%hook:program_hook}",
  }
}

// PatternRules requires various parameters.
type PatternRules struct {
  At reader.Position `if:"internal"`
  PatternRule *[]PatternRule
}

func (*PatternRules) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_rules",
    Spec: "{*pattern_rule}",
  }
}

// PatternType swaps between various options
type PatternType struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*PatternType) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_type",
    Spec: "an {activity:patterned_activity} or a {value:variable_type}",
  }
}

func (*PatternType) Choices() map[string]interface{} {
  return map[string]interface{} {
    "activity": (*PatternedActivity)(nil),
    "value": (*VariableType)(nil),
  }
}

// PatternVariablesDecl requires various parameters.
type PatternVariablesDecl struct {
  At reader.Position `if:"internal"`
  PatternName PatternName
  VariableDecl []VariableDecl
}

func (*PatternVariablesDecl) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_variables_decl",
    Desc: `Declare pattern variables: Storage for values used during the execution of a pattern.`,
    Spec: "The pattern {pattern_name|quote} requires {+variable_decl|comma-and}.",
  }
}

// PatternVariablesTail requires various parameters.
type PatternVariablesTail struct {
  At reader.Position `if:"internal"`
  VariableDecl []VariableDecl
}

func (*PatternVariablesTail) Compose() composer.Spec {
  return composer.Spec{
    Name: "pattern_variables_tail",
    Desc: `Pattern variables: Storage for values used during the execution of a pattern.`,
    Spec: "It requires {+variable_decl|comma-and}.",
  }
}

// PatternedActivity requires a user-specified string.
type PatternedActivity struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *PatternedActivity) String() string {
  return op.Str
}

func (*PatternedActivity) Choices() (closed bool, choices map[string]string) {
  return true, map[string]string{
    "$ACTIVITY": "activity",
  }
}

func (*PatternedActivity) Compose() composer.Spec {
  return composer.Spec{
    Name: "patterned_activity",
    Spec: "{an activity%activity}",
  }
}

// PluralKinds requires a user-specified string.
type PluralKinds struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *PluralKinds) String() string {
  return op.Str
}

func (*PluralKinds) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*PluralKinds) Compose() composer.Spec {
  return composer.Spec{
    Name: "plural_kinds",
    Desc: `Kinds: The plural name of a type of similar nouns. For example: animals, containers, etc.`,
    Spec: "{plural_kinds}",
  }
}

// PrimitiveFunc swaps between various options
type PrimitiveFunc struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*PrimitiveFunc) Compose() composer.Spec {
  return composer.Spec{
    Name: "primitive_func",
    Spec: "{a number%number_eval}, {some text%text_eval}, {a true/false value%bool_eval}",
  }
}

func (*PrimitiveFunc) Choices() map[string]interface{} {
  return map[string]interface{} {
    "number_eval": (*NumberEval)(nil),
    "text_eval": (*TextEval)(nil),
    "bool_eval": (*BoolEval)(nil),
  }
}

// PrimitivePhrase requires various parameters.
type PrimitivePhrase struct {
  At reader.Position `if:"internal"`
  PrimitiveType PrimitiveType
  Property Property
}

func (*PrimitivePhrase) Compose() composer.Spec {
  return composer.Spec{
    Name: "primitive_phrase",
    Spec: "{primitive_type} called {property}",
  }
}

// PrimitiveType requires a user-specified string.
type PrimitiveType struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *PrimitiveType) String() string {
  return op.Str
}

func (*PrimitiveType) Choices() (closed bool, choices map[string]string) {
  return true, map[string]string{
    "$NUMBER": "number", "$TEXT": "text", "$BOOL": "bool",
  }
}

func (*PrimitiveType) Compose() composer.Spec {
  return composer.Spec{
    Name: "primitive_type",
    Spec: "{a number%number}, {some text%text}, or {a true/false value%bool}",
  }
}

// PrimitiveValue swaps between various options
type PrimitiveValue struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*PrimitiveValue) Compose() composer.Spec {
  return composer.Spec{
    Name: "primitive_value",
    Spec: "{text%boxed_text} or {number%boxed_number}",
  }
}

func (*PrimitiveValue) Choices() map[string]interface{} {
  return map[string]interface{} {
    "boxed_text": (*BoxedText)(nil),
    "boxed_number": (*BoxedNumber)(nil),
  }
}

// ProgramHook swaps between various options
type ProgramHook struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*ProgramHook) Compose() composer.Spec {
  return composer.Spec{
    Name: "program_hook",
    Spec: "run an {activity} or return a {result:program_return}",
  }
}

func (*ProgramHook) Choices() map[string]interface{} {
  return map[string]interface{} {
    "activity": (*Activity)(nil),
    "result": (*ProgramReturn)(nil),
  }
}

// ProgramResult swaps between various options
type ProgramResult struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*ProgramResult) Compose() composer.Spec {
  return composer.Spec{
    Name: "program_result",
    Spec: "a {simple value%primitive:primitive_func} or an {object:object_func}",
  }
}

func (*ProgramResult) Choices() map[string]interface{} {
  return map[string]interface{} {
    "primitive": (*PrimitiveFunc)(nil),
    "object": (*ObjectFunc)(nil),
  }
}

// ProgramReturn requires various parameters.
type ProgramReturn struct {
  At reader.Position `if:"internal"`
  Result ProgramResult
}

func (*ProgramReturn) Compose() composer.Spec {
  return composer.Spec{
    Name: "program_return",
    Spec: "return {result:program_result}",
  }
}

// Pronoun requires a user-specified string.
type Pronoun struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Pronoun) String() string {
  return op.Str
}

func (*Pronoun) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    "$IT": "it", "$THEY": "they",
  }
}

func (*Pronoun) Compose() composer.Spec {
  return composer.Spec{
    Name: "pronoun",
    Spec: "{it}, {they}, or {pronoun}",
  }
}

// Property requires a user-specified string.
type Property struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Property) String() string {
  return op.Str
}

func (*Property) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*Property) Compose() composer.Spec {
  return composer.Spec{
    Name: "property",
    Spec: "{property}",
  }
}

// PropertyPhrase swaps between various options
type PropertyPhrase struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*PropertyPhrase) Compose() composer.Spec {
  return composer.Spec{
    Name: "property_phrase",
    Spec: "{primitive_phrase} or {aspect_phrase}",
  }
}

func (*PropertyPhrase) Choices() map[string]interface{} {
  return map[string]interface{} {
    "primitive_phrase": (*PrimitivePhrase)(nil),
    "aspect_phrase": (*AspectPhrase)(nil),
  }
}

// RecordList requires various parameters.
type RecordList struct {
  At reader.Position `if:"internal"`
  Kind SingularKind
}

func (*RecordList) Compose() composer.Spec {
  return composer.Spec{
    Name: "record_list",
    Spec: "a list of {kind%kind:singular_kind} records",
  }
}

// RecordType requires various parameters.
type RecordType struct {
  At reader.Position `if:"internal"`
  Kind SingularKind
}

func (*RecordType) Compose() composer.Spec {
  return composer.Spec{
    Name: "record_type",
    Spec: "a record of {kind%kind:singular_kind}",
  }
}

// Relation requires a user-specified string.
type Relation struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Relation) String() string {
  return op.Str
}

func (*Relation) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*Relation) Compose() composer.Spec {
  return composer.Spec{
    Name: "relation",
    Spec: "{relation}",
  }
}

// RelativeToNoun requires various parameters.
type RelativeToNoun struct {
  At reader.Position `if:"internal"`
  Relation Relation
  Nouns []NamedNoun
  AreBeing AreBeing
  Nouns1 []NamedNoun
}

func (*RelativeToNoun) Compose() composer.Spec {
  return composer.Spec{
    Name: "relative_to_noun",
    Spec: "{relation} {nouns+named_noun} {are_being} {nouns+named_noun}.",
  }
}

// RenderTemplate requires various parameters.
type RenderTemplate struct {
  At reader.Position `if:"internal"`
  Template Lines
}

func (*RenderTemplate) Compose() composer.Spec {
  return composer.Spec{
    Name: "render_template",
    Desc: `Render template: Parse text using iffy templates. See: https://github.com/ionous/iffy/wiki/Templates`,
    Group: "format",
  }
}

// ShuffleText requires various parameters.
type ShuffleText struct {
  At reader.Position `if:"internal"`
  Parts []rt.TextEval
}

func (*ShuffleText) Compose() composer.Spec {
  return composer.Spec{
    Name: "shuffle_text",
    Desc: `Shuffle text: When called multiple times returns its inputs at random.`,
    Group: "format",
  }
}

// SingularKind requires a user-specified string.
type SingularKind struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *SingularKind) String() string {
  return op.Str
}

func (*SingularKind) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*SingularKind) Compose() composer.Spec {
  return composer.Spec{
    Name: "singular_kind",
    Desc: `Kind: Describes a type of similar nouns. For example: an animal, a container, etc.`,
    Spec: "{singular_kind}",
  }
}

// StoppingText requires various parameters.
type StoppingText struct {
  At reader.Position `if:"internal"`
  Parts []rt.TextEval
}

func (*StoppingText) Compose() composer.Spec {
  return composer.Spec{
    Name: "stopping_text",
    Desc: `Stopping text: When called multiple times returns each of its inputs in turn, sticking to the last one.`,
    Group: "format",
  }
}

// Story requires various parameters.
type Story struct {
  At reader.Position `if:"internal"`
  Paragraph *[]Paragraph
}

func (*Story) Compose() composer.Spec {
  return composer.Spec{
    Name: "story",
    Spec: "{*paragraph}",
  }
}


// Summary requires various parameters.
type Summary struct {
  At reader.Position `if:"internal"`
  Lines Lines
}

func (*Summary) Compose() composer.Spec {
  return composer.Spec{
    Name: "summary",
    Spec: "The summary is: {summary%lines|quote}",
  }
}

// Tail requires various parameters.
type Tail struct {
  At reader.Position `if:"internal"`
  Pronoun Pronoun
  NounPhrase NounPhrase
}

func (*Tail) Compose() composer.Spec {
  return composer.Spec{
    Name: "tail",
    Desc: `Trailing statement: Adds details about the preceding noun or nouns.`,
    Spec: "{pronoun} {noun_phrase}.",
  }
}

// TestName requires a user-specified string.
type TestName struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *TestName) String() string {
  return op.Str
}

func (*TestName) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    "$CURRENT_TEST": "current_test",
  }
}

func (*TestName) Compose() composer.Spec {
  return composer.Spec{
    Name: "test_name",
    Spec: "{the test%current_test}, or {test name%test_name}",
  }
}

// TestOutput requires various parameters.
type TestOutput struct {
  At reader.Position `if:"internal"`
  Lines Lines
}

func (*TestOutput) Compose() composer.Spec {
  return composer.Spec{
    Name: "test_output",
    Desc: `Test output: Expect that a test uses 'Say' to print some specific text.`,
    Spec: "output {lines|quote}.",
  }
}

// TestRule requires various parameters.
type TestRule struct {
  At reader.Position `if:"internal"`
  TestName TestName
  Hook ProgramHook
}

func (*TestRule) Compose() composer.Spec {
  return composer.Spec{
    Name: "test_rule",
    Spec: "To test {test_name}: {do%hook:program_hook}",
  }
}

// TestScene requires various parameters.
type TestScene struct {
  At reader.Position `if:"internal"`
  TestName TestName
  Story Story
}

func (*TestScene) Compose() composer.Spec {
  return composer.Spec{
    Name: "test_scene",
    Spec: "While testing {test_name}: {story}",
  }
}

// TestStatement requires various parameters.
type TestStatement struct {
  At reader.Position `if:"internal"`
  TestName TestName
  Test Testing
}

func (*TestStatement) Compose() composer.Spec {
  return composer.Spec{
    Name: "test_statement",
    Spec: "Expect {test_name} to {expectation%test:testing}",
  }
}


// Text requires a user-specified string.
type Text struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Text) String() string {
  return op.Str
}

func (*Text) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    "$EMPTY": "empty",
  }
}

func (*Text) Compose() composer.Spec {
  return composer.Spec{
    Name: "text",
    Desc: `Text: A sequence of characters of any length, all on one line. Examples include letters, words, or short sentences.
Text is generally something displayed to the player.
See also: lines.`,
    Spec: "{text} or {empty}",
  }
}

// TextList requires a user-specified string.
type TextList struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *TextList) String() string {
  return op.Str
}

func (*TextList) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*TextList) Compose() composer.Spec {
  return composer.Spec{
    Name: "text_list",
    Spec: "{a list of text%text_list}",
  }
}

// TextValue requires various parameters.
type TextValue struct {
  At reader.Position `if:"internal"`
  Text Text
}

func (*TextValue) Compose() composer.Spec {
  return composer.Spec{
    Name: "text_value",
    Desc: `Text value: specify a small bit of text.`,
    Group: "literals",
  }
}

// Trait requires a user-specified string.
type Trait struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *Trait) String() string {
  return op.Str
}

func (*Trait) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*Trait) Compose() composer.Spec {
  return composer.Spec{
    Name: "trait",
    Spec: "{trait}",
  }
}

// TraitPhrase requires various parameters.
type TraitPhrase struct {
  At reader.Position `if:"internal"`
  AreEither AreEither
  Trait []Trait
}

func (*TraitPhrase) Compose() composer.Spec {
  return composer.Spec{
    Name: "trait_phrase",
    Spec: "{are_either} {trait+trait|comma-or}.",
  }
}

// VariableDecl requires various parameters.
type VariableDecl struct {
  At reader.Position `if:"internal"`
  Type VariableType
  Name VariableName
}

func (*VariableDecl) Compose() composer.Spec {
  return composer.Spec{
    Name: "variable_decl",
    Spec: "{type:variable_type} called {name:variable_name}",
  }
}

// VariableName requires a user-specified string.
type VariableName struct {
  At  reader.Position `if:"internal"`
  Str string
}

func (op *VariableName) String() string {
  return op.Str
}

func (*VariableName) Choices() (closed bool, choices map[string]string) {
  return false, map[string]string{
    
  }
}

func (*VariableName) Compose() composer.Spec {
  return composer.Spec{
    Name: "variable_name",
    Spec: "{variable_name}",
  }
}

// VariableType swaps between various options
type VariableType struct {
  At  reader.Position `if:"internal"`
  Opt interface{}
}

func (*VariableType) Compose() composer.Spec {
  return composer.Spec{
    Name: "variable_type",
    Spec: "a {simple value%primitive:primitive_type}, an {object:object_type}, or {other value%ext:ext_type}",
  }
}

func (*VariableType) Choices() map[string]interface{} {
  return map[string]interface{} {
    "primitive": (*PrimitiveType)(nil),
    "object": (*ObjectType)(nil),
    "ext": (*ExtType)(nil),
  }
}

var Slots = []composer.Slot{
  {
    Name: "story_statement",
    Type: (*StoryStatement)(nil),
    Desc: "Phrase",
  },
  {
    Name: "testing",
    Type: (*Testing)(nil),
    Desc: "Run a series of tests.",
  },
}

var Model = []composer.Slat{
  (*An)(nil),
  (*AreAn)(nil),
  (*AreBeing)(nil),
  (*AreEither)(nil),
  (*Argument)(nil),
  (*Arguments)(nil),
  (*Aspect)(nil),
  (*AspectPhrase)(nil),
  (*AspectTraits)(nil),
  (*Bool)(nil),
  (*BoxedNumber)(nil),
  (*BoxedText)(nil),
  (*Certainties)(nil),
  (*Certainty)(nil),
  (*Comment)(nil),
  (*Comments)(nil),
  (*CycleText)(nil),
  (*DetermineAct)(nil),
  (*DetermineBool)(nil),
  (*DetermineNum)(nil),
  (*DetermineNumList)(nil),
  (*DetermineText)(nil),
  (*DetermineTextList)(nil),
  (*Determiner)(nil),
  (*ExtType)(nil),
  (*KindOfNoun)(nil),
  (*KindsOfAspect)(nil),
  (*KindsOfKind)(nil),
  (*KindsPossessProperties)(nil),
  (*Lede)(nil),
  (*Lines)(nil),
  (*LocalDecl)(nil),
  (*NamedNoun)(nil),
  (*NounAssignment)(nil),
  (*NounName)(nil),
  (*NounPhrase)(nil),
  (*NounRelation)(nil),
  (*NounStatement)(nil),
  (*NounTraits)(nil),
  (*NounType)(nil),
  (*Number)(nil),
  (*NumberList)(nil),
  (*ObjectFunc)(nil),
  (*ObjectType)(nil),
  (*OptionalProperty)(nil),
  (*Paragraph)(nil),
  (*PatternActions)(nil),
  (*PatternDecl)(nil),
  (*PatternLocals)(nil),
  (*PatternName)(nil),
  (*PatternRule)(nil),
  (*PatternRules)(nil),
  (*PatternType)(nil),
  (*PatternVariablesDecl)(nil),
  (*PatternVariablesTail)(nil),
  (*PatternedActivity)(nil),
  (*PluralKinds)(nil),
  (*PrimitiveFunc)(nil),
  (*PrimitivePhrase)(nil),
  (*PrimitiveType)(nil),
  (*PrimitiveValue)(nil),
  (*ProgramHook)(nil),
  (*ProgramResult)(nil),
  (*ProgramReturn)(nil),
  (*Pronoun)(nil),
  (*Property)(nil),
  (*PropertyPhrase)(nil),
  (*RecordList)(nil),
  (*RecordType)(nil),
  (*Relation)(nil),
  (*RelativeToNoun)(nil),
  (*RenderTemplate)(nil),
  (*ShuffleText)(nil),
  (*SingularKind)(nil),
  (*StoppingText)(nil),
  (*Story)(nil),
  (*Summary)(nil),
  (*Tail)(nil),
  (*TestName)(nil),
  (*TestOutput)(nil),
  (*TestRule)(nil),
  (*TestScene)(nil),
  (*TestStatement)(nil),
  (*Text)(nil),
  (*TextList)(nil),
  (*TextValue)(nil),
  (*Trait)(nil),
  (*TraitPhrase)(nil),
  (*VariableDecl)(nil),
  (*VariableName)(nil),
  (*VariableType)(nil),
}
