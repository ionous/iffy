class NodeTest {
  constructor(rootItem) {
    this.all= {};
    this.nodes= new Nodes( this.all, "" );
    this.redux= new Redux({
      set(tgt, field, value) {
        tgt[field]= value;
      },
      "delete": (tgt, field) => {
        delete tgt[field];
      },
    }, this.nodes, 100);
    if (rootItem){
      this.nodes.unroll(rootItem);
    }
    this.rootItem= rootItem
  }
  newMutation(node) {
    const state= new MutationState(node);
    return new Mutation(this.redux, state);
  }
  expect(src, ids, reason) {
    const have= src.reduce((a,n)=> a + n.id, "");
    console.assert(have === ids, `${reason} unexpected ids ${have}`);
  }
}

function nodeTests() {
  const testStory= {
    "id": "td1",
    "type": "paragraph",
    "value": {
      "$STORY_STATEMENT": [{
        "id": "td0",
        "type": "story_statement",
        "value": {
          "id": "td5",
          "type": "noun_statement",
          "value": {
            "$LEDE": {
              "id": "td4",
              "type": "lede",
              "value": {
                "$NOUN": [{
                  "id": "td2",
                  "type": "noun",
                  "value": {
                     "$COMMON_NOUN": {
                        "id": "td8",
                        "type": "common_noun",
                        "value": {
                          "$DETERMINER": {
                            "id": "td6",
                            "type": "determiner",

                            "value": "$THE"
                          },
                          "$COMMON_NAME": {
                            "id": "td7",
                            "type": "common_name",
                            "value": "box"
                          }}}}}],
                "$NOUN_PHRASE": {
                  "id": "td3",
                  "type": "noun_phrase",
                  "value": {
                    "$KIND_OF_NOUN": {
                      "id": "td11",
                      "type": "kind_of_noun",
                      "value": {
                        "$ARE_AN": {
                          "id": "td9",
                          "type": "are_an",
                          "value": "$ISA"
                        },
                        "$KIND": {
                          "id": "td10",
                          "type": "singular_kind",
                          "value": "container"
                        }}}}}}}}}}]
    }};
  function _runTest(name, testFn, root) {
    const make= new Make(new Types());
    makeLang(make);
    const test= new NodeTest(root);
    console.log("testing", name);
    testFn(test);
  }
  function runTest(name, testFn, root=testStory) {
    try {
      _runTest(name, testFn, root);
    } catch (error) {
      console.error("FAILED", name, error);
    }
  }
  function testMutation(name, expected) {
    runTest(`mutation ${name}`, function(test) {
      const before= JSON.stringify(test.rootItem, 0, 2);
      const node= test.all[name];
      const mutation= test.newMutation(node);
      const have= JSON.stringify(mutation.state,0,2);
      const want= JSON.stringify(expected,0,2);
      if (have !== want) {
        console.log("have:", have);
        console.log("want:", want);
        throw new Error(`${node.id} mismatched`);
      }
      const after= JSON.stringify(test.rootItem, 0, 2);
      if (before !== after) {
        console.log("have:", after);
        throw new Error(`original data changed?!`);
      }
    });
  }

  // "the" ( a determiner ) box
  testMutation("td6", {
    // insert a noun,
    left: [{
      "parent": "td4",
      "token": "$NOUN",
      "item": "td2",
    },{
      "parent": "td1",
      "token": "$STORY_STATEMENT",
      "item": "td0"
    }],
    right: [
    // right of td6 is "$COMMON_NAME", so td6 isnt a right edge.
    ],
    // delete the common noun
    // there's no way to "undo" the common noun choice except by way of delete
    removes: {
      "parent": "td2",
      "token": null,  // choice is $COMMON_NOUN
      "item": "td8"
    }
  });

  // the "box" td7
  testMutation("td7", {
    left: [],
    // right - append a noun
    right: [{
      "parent": "td4",
      "token": "$NOUN",
      "item": "td2",
    }],
    // delete the common noun
    // there's no way to "undo" the common noun choice except by way of delete
    removes: {
      "parent": "td2",
      "token": null, // choice is "$COMMON_NOUN",
      "item": "td8"
    }
  });
  // "is a" container ( td9 is an "are_an" in a noun phrase, in a lede )
  testMutation("td9", {
    left: [{
      // lede is a run: run("lede", "{+noun|comma-and} {noun_phrase}.")
      // to the left of the phrase is a repeating noun slot ( filled a common noun "the box" )
      // so we should be able to add a noun to the end of that array.
      "parent": "td4",
      "token": "$NOUN",
      "item": "td2" // note: previously, null would indicate a terminal addition.
    }],
    right: [{
      "parent": "td11",
      "token": "$TRAIT",
      "item": null // b/c optional
    }],
    removes: { // delete the noun phrase
      "parent": "td3",
      "token": null,
      "item": "td11"
    }
  });

  // "container"
  testMutation("td10", {
    left: [{
      "parent": "td11",
      "token": "$TRAIT",
      "item": null
    }],
    right: [{
      "parent": "td11",
      "token": "$NOUN_RELATION",
      "item": null
    },{
      "parent": "td5",
      "token": "$TAIL",
      "item": null
    },{
      "parent": "td5",
      "token": "$SUMMARY",
      "item": null
    },{
      "parent": "td1",
      "token": "$STORY_STATEMENT",
      "item": "td0"
    }],
    removes: { // delete the noun phrase
      "parent": "td3",
      "token": null, // choice is "$KIND_OF_NOUN"
      "item": "td11"
    }
  });
  //
  runTest("serialization", function(test) {
    const ogJson= JSON.stringify(test.rootItem,0,2);
    const nodeJson= test.nodes.root.serialize();
    if (nodeJson !== ogJson) {
      console.log(nodeJson);
      throw new Error("mismatched serialization");
    }
  });
  //
  // test actually mutating some data
  //
  runTest("test appending to a new (optional) list", function(test) {
    const kindOfNoun= test.all.td11;
    if (kindOfNoun.getKid("$TRAIT")) {
      throw new Error("unexpected initial attribute");
    }
    test.newMutation( test.all.td10 ).mutate(-1);
    if (kindOfNoun.getKid("$TRAIT").length !== 1) {
      throw new Error("missing new attribute");
    }
  });
  //
  runTest("test appending to an existing (required) list", function(test) {
    const lede= test.all.td4;
    if (lede.getKid("$NOUN").length!== 1) {
        throw new Error("expected one initial noun");
    }
    test.newMutation( test.all.td9 ).mutate(-1);
    if (lede.getKid("$NOUN").length!== 2) {
        throw new Error("expected a new noun");
    }
  });
  //
  runTest("test creating a non-repeating optional item", function(test) {
    const nounStatement= test.all.td5;
    if (nounStatement.getKid("$SUMMARY")) {
      throw new Error("unexpected initial summary");
    }
    test.newMutation( test.all.td10 ).mutate(3);
    const summary= nounStatement.getKid("$SUMMARY");
    if (!summary || (summary.type !== "summary")) {
      throw new Error("unexpected new empty summary");
    }
  });
  // deletion,
  runTest("delete a slat", function(test) {
    // td3 is a noun_phrase set to "$KIND_OF_NOUN"
    const nounPhrase= test.all.td3;
    if ((nounPhrase.choice !== "$KIND_OF_NOUN") || (nounPhrase.kid.type != "kind_of_noun")) {
      throw new Error("expected initial 'kind of noun' phrase");
    }
    // delete the "kind of noun" noun phrase
    test.newMutation( test.all.td9 ).mutate(0);
    if (nounPhrase.choice || nounPhrase.kid) {
      throw new Error(`expected noun phrase was deleted ${JSON.stringify(nounPhrase)}`);
    }
  });
  runTest("delete a repeating item", function(test) {
    const statementList= test.all.td1.getKid("$STORY_STATEMENT");
    if (statementList.length !== 2) {
      throw new Error("expected two initial statements");
    }
    // delete the second story statement
    test.newMutation( test.all.td2 ).mutate(0);
    if (statementList.length !== 1) {
      throw new Error("expected one remaining statements");
    }
    if (statementList[0].id !== "td0") {
      throw new Error("expected td0 remains");
    }
  },{
      "id": "td1",
      "type": "story_statements",
      "value": {
        "$STORY_STATEMENT": [{
          "id": "td0",
          "type": "story_statement"
        },{
          "id": "td2",
          "type": "story_statement"
        }]
      }
    });
  runTest("delete an optional item", function(test) {
    const nounStatement= test.all.td1;
    if (!nounStatement.getKid("$SUMMARY")) {
      throw new Error("expected summary statement");
    }
    // delete the summary
    test.newMutation( test.all.td3 ).mutate(0);
    if (nounStatement.getKid("$SUMMARY")) {
      throw new Error("expected no summary statement");
    }
  }, {
    "id": "td0",
    "type": "story_statement",
    "value": {
      "id": "td1",
      "type": "noun_statement",
      "value": {
        "$LEDE": {
          "id": "td2",
          "type": "lede",
          "value": {}
        },
        "$SUMMARY": {
          "id": "td3",
          "type": "summary",
          "value": ""
        }
      }
    }
  });
  //
  runTest("add to left side of root", function(test) {
    const statementList= test.all.td1.getKid("$STORY_STATEMENT");
    if (statementList.length !== 1) {
      throw new Error("expected one additional statement");
    }
    const og= statementList[0];
    test.newMutation( test.all.td6 ).mutate(-2);
    if (statementList.length !== 2) {
      throw new Error("expected one additional statement");
    }
    const ogLeftMost=((statementList[0] === og) &&
                      (statementList[1] !== og));
    // we're adding to the left, so og shouldn't be on the left.
    if (ogLeftMost) {
      throw new Error("expected right side addition");
    }
  });
  runTest("add to right side of root", function(test){
    const statementList= test.all.td1.getKid("$STORY_STATEMENT");
    if (statementList.length !== 1) {
      throw new Error("expected one initial statements");
    }
    const og= statementList[0];
    test.newMutation( test.all.td10 ).mutate(4);
    if (statementList.length !== 2) {
      throw new Error("expected one additional statement");
    }
    const ogLeftMost=((statementList[0] === og) &&
                      (statementList[1] !== og));
    // we're adding to the right, so og should still be on the left.
    if (!ogLeftMost) {
      throw new Error("expected right side addition");
    }
  });
  //
  runTest("add blank story statement", function(test) {
    const para= test.nodes.root;
    const statements= para.getKid("$STORY_STATEMENT");
    const table= new StatementTable(test.redux, para);
    if (statements.length!==1 || statements[0].id !== "td0") {
      throw new Error(" td0 should start as the first statement");
    }
    table.addBlank(0);
    if (statements.length!==2 || statements[1].id !== "td0") {
      throw new Error(" td0 should now be the second statement");
    }
    const isBlank= statements[0];
    if (isBlank.type !== "story_statement" || isBlank.kid !== null) {
      throw new Error("a blank statement should now lead");
    }
    const undid= test.redux.undo();
    if (!undid) {
      throw new Error("undo failed");
    }
    if (statements.length!==1 || statements[0].id !== "td0") {
      throw new Error(" undo should restore td0 as the first statement");
    }
  });
  const emptyParagraph= {
      "id": "root",
      "type": "paragraph",
      "value": {"$STORY_STATEMENT": []}
  };
  runTest("move story statements", function(test) {
    const para= test.nodes.root;
    const statements= para.getKid("$STORY_STATEMENT");
    const table= new StatementTable(test.redux, para);
    table.addBlank(0);
    table.addBlank(1);
    table.addBlank(2);
    test.expect(statements, "012", "initially");
    table.move(1,0,2);
    test.expect(statements, "120", "moved src>dst");
    test.redux.undo();
    test.expect(statements, "012", "undone");
    table.move(0,3,2);
    test.expect(statements, "201", "moved dst>src");
    test.redux.undo();
    test.expect(statements, "012", "undone again");
    const nothrow= true;
    const illegalMove= table.move(0,1,3, nothrow);
    console.assert(illegalMove, "expected illegal move detected")
    test.expect(statements, "012", "steady state");
    table.move(2,0,10000);
    test.expect(statements, "201", "width cap");
    test.redux.undo();
    test.expect(statements, "012", "undone done");
  }, emptyParagraph);
  //
  runTest("drop p from ps to other ps", function(test) {
    // there's only one list of paragraphs per story, so this never happens in reality.
    const { nodes, redux } = test;
    const mainStory= nodes.newFromType(null, "story");
    const otherStory= nodes.newFromType(null, "story");

    const ps1= new ParagraphTable(redux, mainStory);
    const ps2= new ParagraphTable(redux, otherStory);
    test.expect(ps1.items, "1");
    test.expect(ps2.items, "4");

    ps2.transferTo(0, { list:ps1, idx:0 });
    test.expect(ps1.items, "");
    test.expect(ps2.items, "14");

    test.redux.undo();
    test.expect(ps1.items, "1");
    test.expect(ps2.items, "4");
  }, null);
  runTest("drop p from ps appending to a line", function(test) {
  });
  runTest("drop partial line from p into ps, creating a p", function(test) {
  });
  runTest("drop partial line from p appending to a line", function(test) {
  });
  runTest("drop full line from p appending to a line, removing the original line", function(test) {
  });
  runTest("drop full line from p into ps, creating a p, removing the original line", function(test) {
  });


}
nodeTests();
