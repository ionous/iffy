function localLang(make) {
  make.group("Story Statements", function() {
    make.run("story_statements", "{+story_statement}"); // old

    make.run("story", "{+paragraph}");
    make.run("paragraph", "{+story_statement}");
    make.slot("story_statement");
    //
    make.run("noun_statement", "story_statement", "{:lede} {*tail} {?summary}",
             "Noun statement: Describes people, places, or things.");

    make.run("comment", "story_statement", "Note: {comment%lines}",
      "Comment: Information about the story not used by the story.")
  });

  make.group("Tests", function() {
    // "testing" is an interface, currently with one implementation type: TestOutput
    make.run("test_statement", "story_statement",
      "The test {test name%name:text|quote} {expects%test:testing}");

    make.run("test_scene", "story_statement",
      "During the test {test name%name:text|quote}: {+story_statement}");
  });

  make.group("Nouns", function() {
    make.run("lede", "{+noun|comma-and} {noun_phrase}.",
              "Leading statement: Describes one or more nouns.");

    make.run("tail", "{pronoun} {noun_phrase}.",
             "Trailing statement: Adds details about the preceding noun or nouns.");

    // fix? change this into some sort of default pick of noun assignment.
    // make.run("summary", "{The [summary] is:: %lines}");
    make.run("summary", "The summary is: {summary%lines|quote}");

    make.opt("noun_phrase", "{kind_of_noun}, {noun_traits}, or {noun_relation}");

    // fix: think this should always be "are" never "is"
    make.run("kind_of_noun", "{are_an} {trait*trait|comma-and} {kind:singular_kind} {?noun_relation}");

    make.opt("noun", "{proper_noun} or {common_noun}");
    make.run("common_noun", "object_ref", "{determiner} {name:common_name}");
    make.run("proper_noun", "object_ref", "{name:proper_name}");

    make.run("noun_type",  "{an} {kind of%kinds:plural_kinds} noun");

    make.str("proper_name", `Proper Name: A name given to some specific person, place, or thing.
Proper names are usually capitalized. For example, maybe: 'Haruki', 'Jane', or 'Toronto'.`);

    make.str("common_name", `Common Name: A generalized name given to some specific item, place, or thing.
    Common names are usually not capitalized. For example, maybe: 'table', 'chair', or 'dog park'.`);
  });

  make.group("Patterns", function() {
    make.run("pattern_decl", "story_statement",
       "The pattern {name:pattern_name|quote} determines {type:pattern_type}. {optvars?pattern_variables_tail} {about?comments}",
       `Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it.
  Each function in a given pattern has "guards" which determine whether the function applies in a particular situtation.`
     );

    make.run("comments", "{lines|quote}");

    make.run("pattern_variables_decl", "story_statement",
      "The pattern {pattern_name|quote} requires {+variable_decl|comma-and}.",
       `Declare pattern variables: Storage for values used during the execution of a pattern.`);

    make.run("pattern_variables_tail", "It requires {+variable_decl|comma-and}.",
       `Pattern variables: Storage for values used during the execution of a pattern.`);

    make.opt("pattern_type", "an {activity:patterned_activity} or a {value:variable_type}");
    make.str("patterned_activity", "{an activity%activity}");
    make.str("pattern_name");

    // pattern handler
    // similar to pattern_type, but with statements hooks instead of type declarations
    make.run("pattern_handler", "story_statement",
      "To {name:pattern_name}{filters?pattern_filters}: {hook:pattern_hook}",
      "Pattern Handler: Actions to take when a pattern gets used."
      );

    make.run("pattern_filters", " when {filter+bool_eval}");

    make.opt("pattern_hook", "run an {activity:pattern_activity} or return a {result:pattern_return}");

    // fix? pattern_activity and pattern_return both exist for the sake of appearance only.
    make.run("pattern_activity", "run: {activity%go+execute|ghost}");
    make.run("pattern_return", "return {result:pattern_result}");

    make.opt("pattern_result", "a {simple value%primitive:primitive_func} or an {object:object_func}");
    make.opt("primitive_func", "{a number%number_eval}, {some text%text_eval}, {a true/false value%bool_eval}");
    make.run("object_func", "an object named {name%text_eval}");
  });

  make.group("Relations", function() {
    make.run("noun_relation",  "{?are_being} {relation} {+noun|comma-and}");

    make.run("relative_to_noun", "story_statement",
            "{relation} {+noun} {are_being} {+noun}.");

    make.str("relation");
  });

  make.group("Kinds", function() {

    make.run("kinds_of_kind", "story_statement",
         "{plural_kinds} are a kind of {singular_kind}.");

    make.run("kinds_possess_properties", "story_statement",
              "{plural_kinds} have {determiner} {property_phrase}.");

    make.str("singular_kind",
      `Kind: Describes a type of similar nouns.
For example: an animal, a container, etc.`);

    make.str("plural_kinds",
      `Kinds: The plural name of a type of similar nouns.
For example: animals, containers, etc.`);

  });

  make.group("Variables", function() {
    make.run("variable_decl", "{type:variable_type} called {name:variable_name|quote}");
    make.str("variable_name");

    make.opt("variable_type", "a {simple value%primitive:primitive_type} or an {object:object_type}");
    make.run("object_type",  "{an} {kind of%kind:singular_kind}");
  });

  make.group("Traits", function() {
    make.run("kinds_of_aspect", "story_statement", "{aspect} is a kind of value.");
    make.run("aspect_traits", "story_statement", "{aspect} {trait_phrase}");
    make.run("trait_phrase", "{are_either} {trait+trait|comma-or}.");

    make.run("noun_traits", "{are_being} {trait+trait|comma-and}");
    make.str("aspect");
    make.str("trait");
  });

  make.group("Properties", function() {
    // ex. The description of the nets is xxx
    make.run("noun_assignment", "story_statement",
            // "The {property} of {+noun} is the {[text]:: %lines}",
            "The {property} of {+noun} is {the text%lines|summary}",
            "Noun Assignment: Assign text. Gives a noun one or more lines of text.");

    make.opt("property_phrase", "{primitive_phrase} or {aspect_phrase}");
    make.run("aspect_phrase", "{aspect} {?optional_property}");
    make.run("optional_property", "called {property}");

    make.run("certainties", "story_statement",
              "{plural_kinds} {are_being} {certainty} {trait}.");

    make.str("are_either", "{can be%canbe} {are either%either}");

    make.str("certainty",  "{usually}, {always}, {seldom}, or {never}",
             "Certainty: Whether an trait applies to a kind of noun.");

    make.str("property");
  });

  // primitive types
  make.group("Primitive Types", function() {
    make.run("primitive_phrase", "{primitive_type} called {property}");
    make.str("primitive_type", "{a number%number}, {some text%text}, or {a true/false value%bool}");
    make.opt("primitive_value", "{text%boxed_text} or {number%boxed_number}");

    make.run("boxed_text", "{text}");
    make.run("boxed_number", "{number}");
    make.run("boxed_boolean", "{bool}");

    make.str("bool", "{true} or {false}");
    make.str("text", `A sequence of characters of any length, all on one line.
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
    make.str("determiner", "{a}, {an}, {the}, or {other determiner%determiner}");
    make.str("pronoun",  "{it}, {they}, or {pronoun}");
    make.str("an", "{a} or {an}");
    make.str("are_being",  "{are} or {is}");
    make.str("are_an",  "{are}, {are a%area}, {are an%arean}, {is}, {is a%isa}, {is an%isan}");
  });
}

function makeLang(make) {
  make.group("Code", function() {
    // read spec.js
    for (const t of spec) {
      // where's the right place for this...?
      const d= Object.assign(t,
        t.desc&&{desc:Make.makeDesc(t.name, t.desc)},
      );
      if (d.spec) {
        const tags= TagParser.parse(d.spec);
        d["with"].tokens= tags.keys;
        d["with"].params= tags.args;
      }
      make.types.newType(d);
    }
  });
  make.group("Model", function() {
    localLang(make);
  });
}
