class Nodes {
    constructor(rootItem) {
      this.all= {};
      this.makeNodes(rootItem);
      this.redux= new Redux(null, 100);
    }
    newMutation(node) {
      const state= new MutationState(node);
      return new Mutation(this.redux, state);
    }
    // recursively create nodes from data.
    makeNodes( item, parent=null, fieldInParent= null ) {
      const node= new Node( item, parent, fieldInParent );
      this.all[item.id]= node;
      const { value }= item;
      const valueType= typeof value;
      if (valueType === 'string' || valueType === 'number') {
      } else {
        // ex. story statement doesnt have any fields
        // its value is directly another item -- a slat.
        if (value && "id" in value) {
          this.makeNodes( value, node );
        } else {
          // walk fields of our
          for (const field in value) {
            const arg= value[field];
            if (!Array.isArray(arg)) {
              this.makeNodes( arg, node, field );
            } else {
              for (const el of arg) {
                this.makeNodes( el, node, field );
              }
            }
          }
        }
      }
      return node;
    }
}

function nodeTests() {
  function runTest(name, testFn, root) {
    const make= new Make(new Types());
    makeLang(make);

    const nodes= new Nodes(root || {
      "id": "td1",
      "type": "story_statements",
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
                        }
                      }
                    }
                  }],
                  "$NOUN_PHRASE": {
                    "id": "td3",
                    "type": "noun_phrase",
                    "value": {
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
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }]
      }
    });
    console.log("testing", name);
    testFn(nodes);
  }

  function testMutation(name, expected) {
    runTest(`mutation ${name}`, function(nodes) {
      const node= nodes.all[name];
      const mutation= nodes.newMutation(node);
      const have= JSON.stringify(mutation.state,0,2);
      const want= JSON.stringify(expected,0,2);
      if (have !== want) {
        console.log("have:", have);
        console.log("want:", want);
        throw new Error(`${node.item.id} mismatched`);
      }
    });
  }

  // "the" td6
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
    right: [],
    // delete the common noun
    // there's no way to "undo" the common noun choice except by way of delete
    removes: {
      "parent": "td2",
      "token": null,
      "item": "td8"
    }
  });

  // "box" td7
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
      "token": null,
      "item": "td8"
    }
  });
  testMutation("td9", {
    left: [{
      "parent": "td4",
      "token": "$NOUN",
      // "item": "td8"   // in a repeat, you'll want to know what before or after
                         // ( and we'll know which side b/c left or right )
                         // but, when there is no item: we know its a terminal addition
                         // left goes on right side, right goes on left side
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

  // "is a"
  testMutation("td9", {
    left: [{
      "parent": "td4",
      "token": "$NOUN",
      // "item": "td8"   // in a repeat, you'll want to know what before or after
                         // ( and we'll know which side b/c left or right )
                         // but, when there is no item: we know its a terminal addition
                         // left goes on right side, right goes on left side
    }],
    right: [{
        "parent": "td11",
        "token": "$TRAIT",
        "item": null // b/c optional
      }
    ],
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
      "token": null,
      "item": "td11"
    }
  });
  //
  // test actually mutating some data
  //
  runTest("test appending to a new (optional) list", function(nodes) {
    const kindOfNoun= nodes.all.td11.item;
    if (kindOfNoun.value["$TRAIT"]) {
      throw new Error("unexpected initial attribute");
    }
    nodes.newMutation( nodes.all.td10 ).mutate(-1);
    if (kindOfNoun.value["$TRAIT"].length !== 1) {
      throw new Error("missing new attribute");
    }
  });
  //
  runTest("test appending to an existing (required) list", function(nodes) {
    const lede= nodes.all.td4.item;
    if (lede.value["$NOUN"].length !== 1) {
        throw new Error("expected one initial noun");
    }
    nodes.newMutation( nodes.all.td9 ).mutate(-1);
    if (lede.value["$NOUN"].length !== 2) {
        throw new Error("expected a new noun");
    }
  });
  //
  runTest("test creating a non-repeating optional item", function(nodes) {
    const nounStatement= nodes.all.td5.item;
    if (nounStatement.value["$SUMMARY"]) {
      throw new Error("unexpected initial summary");
    }
    nodes.newMutation( nodes.all.td10 ).mutate(3);
    const summary= nounStatement.value["$SUMMARY"];
    if (!summary || (summary.type !== "summary")) {
      throw new Error("unexpected new empty summary");
    }
  });
  // deletion
  runTest("delete a slat", function(nodes) {
    const nounPhrase= nodes.all.td3.item;
    if (nounPhrase.value.type !== "kind_of_noun") {
      throw new Error("expected initial 'kind of noun' phrase");
    }
    // delete the "kind of noun" noun phrase
    nodes.newMutation( nodes.all.td9 ).mutate(0);
    if (nounPhrase.value) {
      throw new Error("expected noun phrase slat was deleted");
    }
  });
  runTest("delete a repeating item", function(nodes) {
    const statementList= nodes.all.td1.item.value["$STORY_STATEMENT"];
    if (statementList.length !== 2) {
      throw new Error("expected two initial statements");
    }
    // delete the second story statement
    nodes.newMutation( nodes.all.td2 ).mutate(0);
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
  runTest("delete an optional item", function(nodes) {
    const nounStatement= nodes.all.td1.item.value;
    if (!nounStatement["$SUMMARY"]) {
      throw new Error("expected summary statement");
    }
    // delete the summary
    nodes.newMutation( nodes.all.td3 ).mutate(0);
    if (nounStatement["$SUMMARY"]) {
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
  runTest("add to left side of root", function(nodes) {
    const statementList= nodes.all.td1.item.value["$STORY_STATEMENT"];
    if (statementList.length !== 1) {
      throw new Error("expected one additional statement");
    }
    const og= statementList[0];
    nodes.newMutation( nodes.all.td6 ).mutate(-2);
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
  runTest("add to right side of root", function(nodes){
    const statementList= nodes.all.td1.item.value["$STORY_STATEMENT"];
    if (statementList.length !== 1) {
      throw new Error("expected one initial statements");
    }
    const og= statementList[0];
    nodes.newMutation( nodes.all.td10 ).mutate(4);
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
}
nodeTests();
