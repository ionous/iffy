function localLang(make) {
  make.group("Story Statements", function() {
    make.run("story_statements", "{+story_statement}"); // old

    make.run("story", "{+paragraph|ghost}");
    make.run("paragraph", "{+story_statement|ghost}");
    make.slot("story_statement");
  });




   make.group("Testing", function() {
    make.run("test", "story_statement",
      "For the test {test_name:text|quote}, expect the output {lines|quote} when running: {go+execute|ghost}.");
  });

  make.group("Story Statements", function() {
    make.run("story_statements", "{+story_statement}"); // old

    make.run("story", "{+paragraph|ghost}");
    make.run("paragraph", "{+story_statement|ghost}");
    make.slot("story_statement");
    //
    make.run("noun_statement", "story_statement", "{:lede} {*tail} {?summary}",
             "Noun statement: Describes people, places, or things.");
  });

  make.group("Nouns", function() {
    make.run("lede", "{+noun|comma-and} {noun_phrase}.",
              "Leading statement: Describes one or more nouns.");

    make.run("tail", "{pronoun} {noun_phrase}.",
             "Trailing statement: Adds details about the preceding noun or nouns.");

    // fix? change this into some sort of default pick of noun assignment.
    // make.run("summary", "{The [summary] is:: %lines}");
    make.run("summary", "The summary is: {summary%lines|quote}");

    make.opt("noun_phrase", "{kind_of_noun}, {noun_attrs}, or {noun_relation}");

    make.opt("noun", "{proper_noun} or {common_noun}");
    make.run("common_noun", "{determiner} {common_name}");
    make.run("proper_noun", "{proper_name}");

    // fix: think this should always be "are" never "is"
    make.run("kinds_of_thing", "story_statement",
             "{plural_kinds} are a kind of {singular_kind}.");

    make.str("proper_name", `Proper Name: A name given to some specific person, place, or thing.
Proper names are usually capitalized. For example, maybe: 'Haruki', 'Jane', or 'Toronto'.`);

    make.str("common_name", `Common Name: A generalized name given to some specific item, place, or thing.
    Common names are usually not capitalized. For example, maybe: 'table', 'chair', or 'dog park'.`);
  });

  make.group("Relations", function() {
    make.run("noun_relation",  "{?are_being} {relation} {+noun|comma-and}");

    make.run("relative_to_noun", "story_statement",
            "{relation} {+noun} {are_being} {+noun}.");

    make.str("relation");
  });

  make.group("Kinds", function() {
    make.run("kind_of_noun", "{are_an} {*attribute|comma-and} {singular_kind} {?noun_relation}");

    make.run("kinds_possess_properties", "story_statement",
              "{plural_kinds} have {determiner} {property_phrase}.");

    make.str("singular_kind",
      `Kind: Describes a type of similar objects.
For example: an animal, a container, etc.`);

    make.str("plural_kinds",
      `Kinds: The plural name of a type of similar objects.
For example: animals, containers, etc.`);
  });

  make.group("Traits", function() {
    make.run("noun_attrs", "{are_being} {+attribute|comma-and}");
    make.str("quality");
    make.str("qualities");
    make.str("attribute");
  });

  make.group("Properties", function() {
    // ex. The description of the nets is xxx
    make.run("noun_assignment", "story_statement",
            // "The {property} of {+noun} is the {[text]:: %lines}",
            "The {property} of {+noun} is {the text%lines|summary}",
            "Assign text. Gives a noun one or more lines of text.");

    make.opt("property_phrase", "{primitive_phrase} or {quality_phrase}");

    make.run("optional_property", "called {property}");

    make.run("kinds_of_quality", "story_statement",
              "{qualities} are a kind of value.");

    make.run("class_attributes", "story_statement",
              "{plural_kinds} {attribute_phrase}");

    make.run("quality_attributes", "story_statement",
              "{qualities} {attribute_phrase}");

    make.run("attribute_phrase", "{are_either} {+attribute|comma-and}.");

    make.run("certainties", "story_statement",
              "{plural_kinds} {are_being} {certainty} {attribute}.");

    make.run("quality_phrase", "{quality} {?optional_property}");

    make.str("are_either", "{can be%canbe} {are either%either}");
    make.str("certainty",  "{usually}, {always}, {seldom}, or {never}",
             "Certainty: Whether an attribute applies to a kind of noun.");

    make.str("property");
  });


  make.group("Helper Types", function() {
    make.str("determiner", "{a}, {an}, {the}, or {other determiner%determiner}");
    make.str("pronoun",  "{it}, {they}, or {pronoun}");
    make.str("an", "{a} or {an}");
    make.str("are_being",  "{are} or {is}");
    make.str("are_an",  "{are}, {are a%area}, {are an%arean}, {is}, {is a%isa}, {is an%isan}");
  });

  // primitive types
  make.group("Primitive Types", function() {

    make.run("primitive_phrase", "{primitive_type} called {property}");
    make.str("primitive_type", "some {text}, a {number}, a {boolean}, a {kind}");
    make.opt("primitive_value", "{text%boxed_text} or {number%boxed_number}");

    make.run("boxed_text", "{text}");
    make.run("boxed_number", "{number}");
    make.run("boxed_boolean", "{bool}");

    make.str("bool", "{true} or {false}");
    make.str("text", `A sequence of characters of any length, all on one line.
Examples include letters, words, or short sentances.
Text is generally something displayed to the player.
See also: lines.`);

    // fix: bracket style links [] for see also?
    make.txt("lines", `A sequence of characters of any length spanning multiple lines.
Paragraphs are a prime example. Generally lines are some piece of the story that will be displayed to the player.
See also: text.`);
    make.num("number");
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
