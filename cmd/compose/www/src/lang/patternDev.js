function makeLang(make) {
  make.group("Story Statements", function() {
    make.run("story_statements", "{+story_statement}"); // old

    make.run("story", "{+paragraph|ghost}");
    make.run("paragraph", "{+story_statement|ghost}");
    make.slot("story_statement");
  });

  make.group("Patterns", function() {
     make.run("pattern_decl", "story_statement",
       "The pattern {name:pattern_name|quote} determines {type:pattern_type}. {optvars?pattern_variables_tail}",
       `Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it.
  Each function in a given pattern has "guards" which determine whether the function applies in a particular situtation.`
     );

    make.run("pattern_variables_decl", "story_statement",
      "The pattern {pattern_name|quote} uses {+variable_decl|comma-and}.",
       `Declare pattern variables: Storage for values used during the execution of a pattern.`);

     make.run("pattern_variables_tail", "It uses {+vars|comma-and}.",
       `Pattern variables: Storage for values used during the execution of a pattern.`);

     make.opt("pattern_type", "an {activity:pattern_activity} or a {value:variable_type}");
     make.str("pattern_activity", "{activity}");
     make.str("pattern_name");
   });

  // primitive types
  make.group("Primitive Types", function() {
    make.run("variable_decl", "{type:variable_type} ( called {name:variable_name|quote} )");
    //
    make.opt("variable_type", "a {simple value%primitive:primitive_type} or an {object:object_type}");

    make.str("primitive_type", "{a number%number}, {some text%text}, or {a true/false value%bool}");
    make.run("object_type",  "{an} {kind of%kinds:plural_kinds} object");

    make.str("singular_kind",
        `Kind: Describes a type of similar objects.
For example: an animal, a container, etc.`);

    make.str("plural_kinds",
        `Kinds: The plural name of a type of similar objects.
For example: animals, containers, etc.`)

    make.str("variable_name");
    make.str("an", "{a} or {an}");
  });
}
