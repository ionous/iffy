[
  {
    "op": "test",
    "path": {
      "parent": "$..[?(@.type=='execute')]"
    },
    "subpatches": [
      {
        "op": "replace",
        "path": {
          "parent": "$..[?@.type=='named_noun']",
          "field": "value"
        },
        "value": "*************************"
      }
    ]
  }
]
