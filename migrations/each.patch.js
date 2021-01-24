[
  {
    "op": "test",
    "path": {
      "parent": "$..[?(@.type=='list_each')]"
    },
    "reason": "select the nodes containing type #choose.",
    "subpatches": [
      {
        "op": "remove",
        "path": {
          "parent": "$.value",
          "field": "$ELSE"
        },
        "reason": "remove #else block, always empty right now"
      },
      {
        "op": "move",
        "from": {
          "parent": "$.value",
          "field": "$GO"
        },
        "path": {
          "parent": "$.value",
          "field": "$DO"
        },
        "reason": "#go now #do"
      },
      {
        "op": "replace",
        "path": {
          "parent": "$.value",
          "field": "$AS"
        },
        "reason": "add #as seletion",
        "value": {
          "type": "list_iterator",
          "value": {
            "type": "as_???",
            "value": {
              "$VAR": {
                "type": "variable_name",
                "value": "???"
              }
            }
          }
        }
      },
      {
        "op": "move",
        "from": {
          "parent": "$.value['$WITH']",
          "field": "value"
        },
        "path": {
          "parent": "$.value['$AS'].value.value['$VAR']",
          "field": "value"
        },
        "reason": "move var name in"
      },
      {
        "op": "remove",
        "path": {
          "parent": "$.value",
          "field": "$WITH"
        },
        "reason": "remove #with"
      }
    ]
  }
]
