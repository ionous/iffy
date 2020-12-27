function localLang(make) {
  make.group("Story Statements", function() {
    make.flow("story", "{*paragraph}");
    make.flow("paragraph", "{*story_statement}", "Phrases");
    make.slot("story_statement", "Phrase");
    //
    make.flow("noun_statement", "story_statement", "{:lede} {*tail} {?summary}",
             "Declare a noun: Describes people, places, or things.");

    make.flow("comment", ["story_statement", "execute"], "Note: {comment%lines}",
      "Add a note: Information about the story for you and other authors.")
  });

  make.group("Tests", function() {
    // "testing" is an interface, currently with once implementation type: TestOutput
    make.slot("testing", "Run a series of tests.");

    make.flow("test_statement", "story_statement",
      "Expect {test_name} to {expectation%test:testing}",
      "Describe test results");

    make.flow("test_scene", "story_statement",
      "While testing {test_name}: {story}",
      "Create a scene for testing");

    make.flow("test_rule", "story_statement",
      "To test {test_name}: {do%hook:program_hook}",
      "Add actions to a test");

    make.flow("test_output", "testing",
      "output {lines|quote}.",
      `Test Output: Expect that a test uses 'Say' to print some specific text.`);

    // would like just "<author's test name>" to be quoted, and not the current_test determiner.
    make.str("test_name", "{the test%current_test}, or {test name%test_name}");
  });

  make.group("Nouns", function() {
    make.flow("lede", "{nouns+named_noun|comma-and} {noun_phrase}.",
              "Leading statement: Describes one or more nouns.");

    make.flow("tail", "{pronoun} {noun_phrase}.",
             "Trailing statement: Adds details about the preceding noun or nouns.");

    // fix? change this into some sort of default pick of noun assignment.
    // make.flow("summary", "{The [summary] is:: %lines}");
    make.flow("summary", "The summary is: {summary%lines|quote}");

    make.swap("noun_phrase", "{kind_of_noun}, {noun_traits}, or {noun_relation}");

    // fix: think this should always be "are" never "is"
    // fix: this shouldnt be "kind of", kind of declares a kind
    // ( but note singular vs. plural nouns phrases here along with are/is )
    // probably should have a switch for singular/ plural -- would be nice if are_an could look ahead and mutate with a custom filter maybe.
    make.flow("kind_of_noun", "{are_an} {*trait|comma-and} kind of {kind:singular_kind} {?noun_relation}");

    make.flow("named_noun", "object_eval", "{determiner} {name:noun_name}");

    make.str("determiner", "{a}, {an}, {the}, {our}, or {other determiner%determiner}",
      `Determiners: modify a word they are associated to designate specificity or, sometimes, a count.
        For instance: "some" fish hooks, "a" pineapple, "75" triangles, "our" Trevor.`  );

    make.str("noun_name",
      `Noun name: Some specific person, place, or thing; or, more rarely, a kind.
        Proper names are usually capitalized:  For example, maybe: 'Haruki', 'Jane', or 'Toronto'.
        Common names are usually not capitalized. For example, maybe: 'table', 'chair', or 'dog park'.
        A set of duplicate object uses their kind. For instance: twelve 'cats'.`);

    make.str("pronoun",  "{it}, {they}, or {pronoun}");

  });

  make.group("Patterns", function() {
    make.flow("pattern_decl", "story_statement",
       "The pattern {name:pattern_name|quote} determines {type:pattern_type}. {optvars?pattern_variables_tail} {about?comment}",
       `Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it.
  Each function in a given pattern has "guards" which determine whether the function applies in a particular situtation.`
     );

    make.flow("pattern_variables_decl", "story_statement",
      "The pattern {pattern_name|quote} requires {+variable_decl|comma-and}.",
       `Add parameters to a pattern: Values provided when starting pattern.`);

    make.flow("pattern_variables_tail", "It requires {+variable_decl|comma-and}.",
       `Pattern variables: Storage for values used during the execution of a pattern.`);

    make.swap("pattern_type", "an {activity:patterned_activity} or a {value:variable_type}");
    make.str("patterned_activity", "{an activity%activity}");
    make.str("pattern_name");

    make.flow("pattern_actions", "story_statement",
      "To {pattern name%name:pattern_name} {?pattern_locals}:{pattern_rules}",
      "Add actions to a pattern: Actions to take when using a pattern.");

    make.flow("pattern_rules", "{*pattern_rule}");
    make.flow("pattern_rule", `When {conditions are met%guard:bool_eval}{ continue%flags?pattern_flags}, then: {do%hook:program_hook}`,
      "Rule");

    make.str("pattern_flags", "{continue before%before}, {continue after%after}, {terminate}");

    make.flow("pattern_locals", " use {+variable_decl|comma-and}",
      "Local: local variables can use the parameters of a pattern to compute temporary values.");

    make.swap("program_hook", "run an {activity} or return a {result:program_return}");

    // fix? activity and program_return both exist for the sake of appearance only.
    make.flow("program_return", "return {result:program_result}");

    make.swap("program_result", "a {simple value%primitive:primitive_func} or an {object:object_func}");

    make.swap("primitive_func", "{a number%number_eval}, {some text%text_eval}, {a true/false value%bool_eval}");
    make.flow("object_func", "an object named {name:text_eval}");
  });

  make.group("Relations", function() {
    make.flow("noun_relation",  "{?are_being} {relation} {nouns+named_noun|comma-and}");

    make.flow("relative_to_noun", "story_statement",
            "{relation} {nouns+named_noun} {are_being} {nouns+named_noun}.",
            "Relate nouns to each other");

    make.str("relation");
  });

  make.group("Kinds", function() {
    make.flow("kinds_of_kind", "story_statement",
         "{plural_kinds} are a kind of {singular_kind}.",
         "Declare a kind");

    make.flow("kinds_possess_properties", "story_statement",
              "{plural_kinds} have {+property_decl|comma-and}.",
              "Add properties to a kind");

    make.str("singular_kind",
      `Kind: Describes a type of similar nouns.
For example: an animal, a container, etc.`);

    make.str("plural_kinds",
      `Kinds: The plural name of a type of similar nouns.
For example: animals, containers, etc.`);
  });

  make.group("Records", function() {
    make.flow("kinds_of_record", "story_statement",
         "{records%record_plural} are a kind of record.",
         "Declare a record");

    make.flow("records_possess_properties", "story_statement",
              "{records%record_plural} have {+property_decl|comma-and}.",
              "Add properties to a record");

    make.str("record_singular",
      `Record: Describes a common set of properties.`);

    make.str("record_plural",
      `Records: The plural name for a record.`);
  });

  make.group("Variables", function() {
    make.flow("variable_decl", "{an:determiner} {name:variable_name} ( {type:variable_type}  {comment?lines} )");
    make.str("variable_name");

    make.swap("variable_type", "a {simple value%primitive:primitive_type}, an {object:object_type}, or {other value%ext:ext_type}");
    make.flow("object_type",  "{an:ana} {kind of%kind:singular_kind}");
  });

  make.group("Traits", function() {
    make.flow("kinds_of_aspect", "story_statement", "{aspect} is a kind of value.",
      "Declare an aspect");
    make.flow("aspect_traits", "story_statement", "{aspect} {trait_phrase}",
      "Add traits to an aspect");
    make.flow("trait_phrase", "{are_either} {+trait|comma-or}.");

    make.flow("noun_traits", "{are_being} {+trait|comma-and}");
    make.str("aspect");
    make.str("trait");
  });

  make.group("Properties", function() {
    // ex. The description of the nets is xxx
    make.flow("noun_assignment", "story_statement",
            // "The {property} of {+noun} is the {[text]:: %lines}",
            "The {property} of {nouns+named_noun} is {the text%lines|summary}",
            "Assign text to a noun: Assign text. Gives a noun one or more lines of text.");

    make.flow("property_decl", "{an:determiner} {property} ( {property_type} {comment?lines} )");
    make.swap("property_type", "an {aspect%property_aspect}, {simple value%primitive:primitive_type}, or {other value%ext:ext_type}");

    make.str("property_aspect", "{an aspect%aspect}");

    make.flow("certainties", "story_statement",
              "{plural_kinds} {are_being} {certainty} {trait}.",
              "Give a kind a trait");

    make.str("are_either", "{can be%canbe} {are either%either}");

    make.str("certainty",  "{usually}, {always}, {seldom}, or {never}",
             "Certainty: Whether an trait applies to a kind of noun.");

    make.str("property");
  });

  // primitive types
  make.group("Values", function() {
    make.str("primitive_type", "{a number%number}, {some text%text}, or {a true/false value%bool}");
    make.swap("primitive_value", "{text%boxed_text} or {number%boxed_number}");

    // a list of numbers, a list of text, a record, or a list of records.
    make.swap("ext_type", "a list of {numbers:number_list}, a list of {text%text_list}, a {record:record_type} or a list of {records:record_list}.")

    make.flow("record_type",  "a record of {kind%kind:record_singular}");
    make.flow("record_list",  "a list of {kind%kind:record_singular} records");

    make.flow("boxed_text", "{text}");
    make.flow("boxed_number", "{number}");

    // constants
    make.str("text_list", "{a list of text%list}");
    make.str("number_list", "{a list of numbers%list}");

    make.str("bool", "{true} or {false}");
    make.str("text", "{text} or {empty}", `A sequence of characters of any length, all on one line.
Examples include letters, words, or short sentences.
Text is generally something displayed to the player.
See also: lines.`);

    // fix: bracket style links [] for see also?
    make.txt("lines", `A sequence of characters of any length spanning multiple lines.
Paragraphs are a prime example. Generally lines are some piece of the story that will be displayed to the player.
See also: text.`);
    make.num("number");
  });

 make.group("Helper Types", function() {
    make.str("ana", "{a} or {an}");
    make.str("are_being",  "{are} or {is}");
    make.str("are_an",  "{are}, {are a%area}, {are an%arean}, {is}, {is a%isa}, {is an%isan}");
  });
}

function makeLang(make) {
  // read spec.js ( generated by iffy/cmd/spec/spec.go )
  make.group("Code", function() {
    spec.forEach((t)=> {
      make.newFromSpec(t);
    });
  });
  // read the local language
  make.group("Model", function() {
    localLang(make);
  });
}
