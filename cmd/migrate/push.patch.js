[
  {
    "patch": "replace",
    "migration": {
      "from": {
        "parent": "$..[?(@.type=='list_push')].value",
        "field": "$AT_EDGE"
      },
      "with": {
        "type": "list_edge",
        "value": "$FALSE"
      }
    }
  },
  {
    "patch": "replace",
    "why": "add 'into', we will get the value of 'list' then delete 'list'",
    "migration": {
      "from": {
        "parent": "$..[?(@.type=='list_push')].value",
        "field": "$INTO"
      },
      "with": {
        "type": "list_target",
        "value": {
          "type": "into_rec_list",
          "value": {
            "$VAR": {
              "type": "variable_name"
            }
          }
        }
      }
    }
  },
  {
    "patch": "copy",
    "migration": {
      "from": {
        "parent": "$..[?(@.type=='list_push')].value['$LIST']",
        "field": "value"
      },
      "to": {
        "parent": "$..[?(@.type=='list_push')].value['$INTO']..['$VAR']",
        "field": "value"
      }
    }
  },
  {
    "patch": "replace",
    "migration": {
      "from": {
        "parent": "$..[?(@.type=='list_push')].value",
        "field": "$LIST"
      },
      "with": null
    }
  },
  {
    "patch": "copy",
    "migration": {
      "from": {
        "parent": "$..[?(@.type=='list_push')].value",
        "field": "$INSERT"
      },
      "to": {
        "parent": "$..[?(@.type=='list_push')].value",
        "field": "$FROM"
      }
    }
  },
  {
    "patch": "replace",
    "migration": {
      "from": {
        "parent": "$..[?(@.type=='list_push')].value",
        "field": "$INSERT"
      },
      "with": null
    }
  },
  {
    "patch": "replace",
    "migration": {
      "from": {
        "parent": "$..[?(@.type=='list_push')]",
        "field": "type"
      },
      "with": "put_edge"
    }
  }
]
