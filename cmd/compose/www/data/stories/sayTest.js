
function getStory() {
  return {
    "id": "id-3",
    "type": "test",
    "value": {
      "$TEST_NAME": {
        "id": "id-0",
        "type": "text",
        "value": "hello, goodbye"
      },
      "$GO": [
        {
          "id": "id-1",
          "type": "execute",
          "value": {
            "id": "id-7",
            "type": "choose",
            "value": {
              "$FALSE": [
                {
                  "id": "id-4",
                  "type": "execute",
                  "value": {
                    "id": "id-15",
                    "type": "say",
                    "value": {
                      "$TEXT": {
                        "id": "id-14",
                        "type": "text_eval",
                        "value": {
                          "id": "id-20",
                          "type": "text_value",
                          "value": {
                            "$TEXT": {
                              "id": "id-19",
                              "type": "text",
                              "value": "goodbye"
                            }
                          }
                        }
                      }
                    }
                  }
                }
              ],
              "$IF": {
                "id": "id-5",
                "type": "bool_eval",
                "value": {
                  "id": "id-9",
                  "type": "bool_value",
                  "value": {
                    "$BOOL": {
                      "id": "id-8",
                      "type": "bool",
                      "value": "$TRUE"
                    }
                  }
                }
              },
              "$TRUE": [
                {
                  "id": "id-6",
                  "type": "execute",
                  "value": {
                    "id": "id-11",
                    "type": "say",
                    "value": {
                      "$TEXT": {
                        "id": "id-10",
                        "type": "text_eval",
                        "value": {
                          "id": "id-13",
                          "type": "text_value",
                          "value": {
                            "$TEXT": {
                              "id": "id-12",
                              "type": "text",
                              "value": "hello"
                            }
                          }
                        }
                      }
                    }
                  }
                }
              ]
            }
          }
        }
      ],
      "$LINES": {
        "id": "id-2",
        "type": "lines",
        "value": "hello"
      }
    }
  }
}
