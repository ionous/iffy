function directiveTests() {
    const make= new Make(new Types());
    function test(a,b) {
      console.log("testing", a.name);
      a= JSON.stringify(a);
      b= JSON.stringify(b);
      if (a!== b) {
        console.log("got: ", a);
        console.log("expected: ", b);
        throw new Error("mismatch");
      }
    }
    test(make.str("determiner", "{an}, {the}, or {other determiner%determiner}"),{
      name: "determiner",
      uses: "str",
      with: {
        tokens: ["$AN", ", ", "$THE", ", or ", "$DETERMINER"],
        params: {
          "$AN": "an",
          "$THE": "the",
          "$DETERMINER": {
            label: "other determiner",
            value: null,
          }
        }
      }
    });
    test(make.flow("root", "{traits}"),{
      name: "root",
      uses: "flow",
      with: {
        tokens: ["$TRAITS"],
        params: {
          "$TRAITS": {
            label: "traits",
            type: "traits",
          }
        }
      }
    });
    test(make.flow("traits", "{one or more traits%TRAIT*trait}",
                  "a list of states describing a noun"), {
      name: "traits",
      desc: "a list of states describing a noun",
      uses: "flow",
      with: {
        tokens: ["$TRAIT"],
        params: {
          "$TRAIT": {
            label: "one or more traits",
            type: "trait",
            optional: true,
            repeats: true
          }
        }
      }
    });
    test(make.str("trait",
      "a state describing a noun"), {
      name: "trait",
      desc: "a state describing a noun",
      uses: "str",
      with: {
        tokens: ["$TRAIT"],
        params: {
          "$TRAIT": {
            label: "trait",
            value: null,
          }
        }
      }
    });

    test(make.txt("multiline", "descriptive text"), {
      name: "multiline",
      desc: "descriptive text",
      uses: "txt",
      with:{
        tokens:["$MULTILINE"],
        params:{
          "$MULTILINE":{
            label:"multiline",
            value: null
          }
        }
      }
    });
    test(make.str("certainty",  "{usually}, {always:5}, {seldom}, or {never}",
      "whether an trait applies to a kind of noun."), {
      name: "certainty",
      desc: "whether an trait applies to a kind of noun.",
      uses: "str",
      with: {
        tokens: ["$USUALLY", ", ", "$ALWAYS", ", ", "$SELDOM", ", or ", "$NEVER"],
        params: {
          "$USUALLY": "usually",
          "$ALWAYS": {
            label: "always",
            value: 5
          },
          "$SELDOM": "seldom",
          "$NEVER": "never"
        }
      }
    });
    test(make.opt("noun_phrase",
      "the {kind:kind_of_noun}, {traits%traits:noun_traits}, or {relationships%rel:noun_relation} of a noun.",
      "characteristics of the preceding noun or nouns"), {
      name: "noun_phrase",
      desc: "characteristics of the preceding noun or nouns",
      uses: "opt",
      with: {
        tokens: ["the ", "$KIND", ", ", "$TRAITS", ", or ", "$REL", " of a noun."],
        params: {
          "$KIND": {
            label: "kind",
            type: "kind_of_noun"
          },
          "$TRAITS": {
            label: "traits",
            type: "noun_traits"
          },
          "$REL": {
            label: "relationships",
            type: "noun_relation"
          },
        }
      }
    });
    test(make.flow("kind_of_noun", "", "the classification of nouns by type"), {
      name: "kind_of_noun",
      desc: "the classification of nouns by type",
      uses: "flow",
      with: {
        tokens: [],
        params: {},
      }
    });
    test(make.flow("noun_traits", "", "the status of a noun"), {
      name: "noun_traits",
      desc: "the status of a noun",
      uses: "flow",
      with: {
        tokens: [],
        params: {},
      }
    });
  test(make.flow("noun_relation", "", "the relation of nouns to other nouns."), {
    name: "noun_relation",
    desc: "the relation of nouns to other nouns.",
    uses: "flow",
    with: {
      tokens: [],
      params: {},
    }
  });
  test(make.slot("story_statement", "Sentences are the primary unit of stories."), {
    name: "story_statement",
    desc: "Sentences are the primary unit of stories.",
    uses: "slot"
  });
  test(make.flow("noun_statement", "story_statement"), {
      name: "noun_statement",
      uses: "flow",
      with: {
        slots: ["story_statement"],
        tokens: [],
        params: {},
      }
    });
}
directiveTests();
