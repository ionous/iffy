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
    test(make.run("root", "{attributes%ATTRS:attributes}"),{
      name: "root",
      uses: "run",
      with: {
        tokens: ["$ATTRS"],
        params: {
          "$ATTRS": {
            label: "attributes",
            type: "attributes",
          }
        }
      }
    });
    test(make.run("attributes", "{one or more attributes%ATTR*attribute}",
                  "a list of states describing a noun"), {
      name: "attributes",
      desc: "a list of states describing a noun",
      uses: "run",
      with: {
        tokens: ["$ATTR"],
        params: {
          "$ATTR": {
            label: "one or more attributes",
            type: "attribute",
            optional: true,
            repeats: true
          }
        }
      }
    });
    test(make.str("attribute",
      "a state describing a noun"), {
      name: "attribute",
      desc: "a state describing a noun",
      uses: "str",
      with: {
        tokens: ["$ATTRIBUTE"],
        params: {
          "$ATTRIBUTE": {
            label: "attribute",
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
      "whether an attribute applies to a kind of noun."), {
      name: "certainty",
      desc: "whether an attribute applies to a kind of noun.",
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
      "the {kind:kind_of_noun}, {attributes%attrs:noun_attrs}, or {relationships%rel:noun_relation} of a noun.",
      "characteristics of the preceding noun or nouns"), {
      name: "noun_phrase",
      desc: "characteristics of the preceding noun or nouns",
      uses: "opt",
      with: {
        tokens: ["the ", "$KIND", ", ", "$ATTRS", ", or ", "$REL", " of a noun."],
        params: {
          "$KIND": {
            label: "kind",
            type: "kind_of_noun"
          },
          "$ATTRS": {
            label: "attributes",
            type: "noun_attrs"
          },
          "$REL": {
            label: "relationships",
            type: "noun_relation"
          },
        }
      }
    });
    test(make.run("kind_of_noun", "", "the classification of nouns by type"), {
      name: "kind_of_noun",
      desc: "the classification of nouns by type",
      uses: "run",
      with: {
        tokens: [],
        params: {},
      }
    });
    test(make.run("noun_attrs", "", "the status of a noun"), {
      name: "noun_attrs",
      desc: "the status of a noun",
      uses: "run",
      with: {
        tokens: [],
        params: {},
      }
    });
  test(make.run("noun_relation", "", "the relation of nouns to other nouns."), {
    name: "noun_relation",
    desc: "the relation of nouns to other nouns.",
    uses: "run",
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
  test(make.run("noun_statement", "story_statement"), {
      name: "noun_statement",
      uses: "run",
      with: {
        slots: ["story_statement"],
        tokens: [],
        params: {},
      }
    });
}
directiveTests();
