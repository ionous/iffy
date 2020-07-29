
Vue.component('mk-opt-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.item.type"
    ><mk-switch
      v-if="hasPicked"
      :node="childNode"
    ></mk-switch
    ><mk-pick-inline v-else
      :node="node"
      @picked="onPick"
    ></mk-pick-inline
  ></span>`,
  computed: {
    hasPicked() {
      // REFACTOR -- this should be based on node kid, not item.
      const { value } = this.node.item;
      return value && !!Object.keys(value).length;
    },
    childNode() {
      return this.node.firstChild;
    },
  },
  methods: {
    onPick(token) {
      const { node } = this;
      const { params } = node.itemType.with;
      if (!token in params) {
        throw new Error(`unknown token picked '${token}'`);
      }
      const param= params[token];
      const childType= param.type || param; // an opt's param can map straight to their type.
      const childItem= Types.createItem(childType);

      this.$root.setPrim(node, { [token]: childItem } );
      this.childNode= this.$root.nodes.newNode(node, childItem, token);
    },
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type:Node,
      required:true,
    }
  }
});

/*
example_type: {
  "name": "primitive_value",
  "uses": "opt",
  "with": {
    "tokens": [
      "$BOXED_TEXT",
      " or ",
      "$BOXED_NUMBER"
    ],
    "params": {
      "$BOXED_TEXT": {
        "label": "text",
        "type": "text"
      },
      "$BOXED_NUMBER": "number",
      }
    }
  }
},
example_data: {
  "type": "primitive_value",
  "value": {
    "$BOXED_TEXT": {
      "type": "boxed_text",
      "value": {
        "$TEXT": {
          "type": "text",
          "value": "5"
        }
      }
    }
  }
}

*/
