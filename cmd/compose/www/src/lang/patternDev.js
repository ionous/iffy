function makeLang(make) {
  make.group("Story Statements", function() {
    make.run("story_statements", "{+story_statement}"); // old

    make.run("story", "{+paragraph|ghost}");
    make.run("paragraph", "{+story_statement|ghost}");
    make.slot("story_statement");
  });

  make.group("Patterns", function() {
     make.run("pattern_decl", "story_statement",
       "The pattern {pattern_name|quote} determines {pattern_type}. {?pattern_variables_tail}",
       `Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it.
  Each function in a given pattern has "guards" which determine whether the function applies in a particular situtation.`
     );

    make.run("pattern_variables_decl", "story_statement",
      "The pattern {pattern_name|quote} uses {+variable_decl|comma-and}.",
       `Declare pattern variables: Storage for values used during the execution of a pattern.`);

     make.run("pattern_variables_tail", "It uses {+variable_decl|comma-and}.",
       `Pattern variables: Storage for values used during the execution of a pattern.`);

     make.opt("pattern_type", "an {activity%pattern_activity} or a {value%variable_type}");
     make.str("pattern_activity", "{activity}");
     make.str("pattern_name");
   });

  // primitive types
  make.group("Primitive Types", function() {
    // fix: support repetition of properties: have a text called X, and a blarg called Y.
    make.run("variable_decl", "{variable_type} ( called {property_name|quote|default=result} )");
    // link text [] is a bit squirrelly. --
    //   {a kind of [object]%singular_kind}
    // the link "object" shows up in the option line
    // but its the rest of the text shows up when selected.
    // its handled by the str control, looking upward.
    make.opt("variable_type", "a {simple value%primitive_type} or {[kind] object%plural_kinds}");
    make.str("primitive_type", "{a number%number}, {some text%text}, or {a true/false value%bool}");
    make.str("singular_kind",
        `Kind: Describes a type of similar objects.
For example: an animal, a container, etc.`);
    make.str("plural_kinds",
        `Kinds: The plural name of a type of similar objects.
For example: animals, containers, etc.`)
    make.str("property_name");

  });
}
