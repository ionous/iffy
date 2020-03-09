//JSON.stringify(app.$data.story,0,2);
function getStory () {
  return {
    "id": "id1",
    "type": "story_statements",
    "type": "summary",
    "value": {
      "$STORY_STATEMENT": [{
        "id": "id0",
        "type": "story_statement",
        "value": {
          "id": "id5",
          "type": "noun_statement",
          "value": {
            "$LEDE": {
              "id": "id4",
              "type": "lede",
              "value": {
                "$NOUN": [{
                  "id": "id2",
                  "type": "noun",
                  "value": {
                    "id": "id8",
                    "type": "common_noun",
                    "value": {
                      "$DETERMINER": {
                        "id": "id6",
                        "type": "determiner",
                        "value": "$THE"
                      },
                      "$COMMON_NAME": {
                        "id": "id7",
                        "type": "common_name",
                        "value": "box"
                      }
                    }
                  }
                }],
                "$NOUN_PHRASE": {
                  "id": "id3",
                  "type": "noun_phrase",
                  "value": {
                    "id": "id11",
                    "type": "kind_of_noun",
                    "value": {
                      "$ARE_AN": {
                        "id": "id9",
                        "type": "are_an",
                        "value": "$ISA"
                      },
                      "$KIND": {
                        "id": "id10",
                        "type": "kind",
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
  };
}
