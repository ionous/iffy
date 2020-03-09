Vue.component('mk-run-ctrl', {
  template:
    `<span
    :class="bemBlock()"
    :data-tag="node.item.type"
    ><mk-run-param
      v-for="param in params"
      :param="param"
      :key="param.name"
      :data-tag="param.type"
      @ghost="onGhost($event)"
    ></mk-run-param
  ></span>`,
  components: {
    // a single parameter, for repeats may contain multiple elements.
    'mk-run-param': {
      data() {
        return {
          mod: this.param.plainText ? "plain" :
            this.param.repeats ? "repeats" : ""
        };
      },
      props: {
        param: Object,
      },
      mixins: [bemMixin("mk-run-ctrl")],
      template:
     `<span
       :class="bemElem('item', mod)"
      >{{param.head}}<template
        v-if="param.plainText"
      >{{param.plainText}} </template
      ><template v-for="el in param.kids"
        ><template v-if="param.commas"
          ><template v-if="el.index && (!el.last || el.index>1)"
          >, </template
          ><template v-if="el.last"
          > {{param.commas}} </template
        ></template
        ><mk-switch
          :key="el.id"
          :node="el.node"
        ></mk-switch
      ></template
      >{{param.tail}}<mk-a-button
        v-if="param.ghost"
        :class="bemElem('ghost')"
        @activate="$emit('ghost', param.token)"
      >{{param.ghost}}</mk-a-button
    ></span>`,
    },
  },
  computed: {
    // to simplify the template, we combine the tokens, params, and arguments.
    // each argument is turned into an array of items for recursive unrolling.
    // by definition, each part of a run is an item that produces a value,
    // even if the item produces a constant value: ex. str, txt, or num.
    // originally, each token was a component, but components require an (extra) root element.
    params() {
      let params = [];
      const { node } = this;
      if (!node.item.value) {
        // in theory, all items should have at least have empty data { {} )
        // tbd how true that is. ex. for test data, bad serialization....
        // when this there is no value, there will be no sub controls.
        const msg = `[WARN] item ${node.item.id} has null data`;
        console.log(msg);
        params.push({
          plainText: msg
        });
      } else {
        const { item, itemType: { with: spec } } = node;
        for (const token of spec.tokens) {
          if (!token.startsWith("$")) {
            params.push({
              plainText: token
            });
          } else {
            const arg = item.value[token];
            // note: not every parameter is filled with an arg.
            // we dont worry about inventing required args here.
            // we leave that up to the deserializer ( createItem )
            if (arg) {
              const param = spec.params[token];
              let commas = false;
              let ghost= false;
              let head="";
              let tail= "";
              if (param.filters) {
                if (param.filters.includes("quote")) {
                  head= `\u201C`;
                  tail= `\u201D`;
                }
                if (param.repeats) {
                  if (arg.length > 1) {
                    if (param.filters.includes("comma-and")) {
                      commas = "and";
                    } else if (param.filters.includes("comma-or")) {
                      commas = "or";
                    }
                  }
                  if (param.filters.includes("ghost")) {
                    const gtype= Types.get(param.type);
                    ghost= Types.labelOf(gtype);
                  }
                }
              }
              // turn all elements into a list of elements.
              const childItems = param.repeats ? arg : [arg];
              const kids = childItems.map((childItem, i) => {
                return {
                  id: childItem.id,
                  node: node.newKid(childItem, token),
                  index: i,
                  last: i === (childItems.length - 1)
                }
              });
              const name = token.replace("$", "").toLowerCase();
              const type = param.type;
              params.push({
                name,
                type,
                head,
                tail,
                kids,
                commas,
                token,
                ghost,
                repeats:param.repeats
              });
            }
          }
        };
      }
      return params;
    }
  },
  methods: {
    // when the ghost is clicked, we want to expand it.
    onGhost(token) {
      this.$root.newGhost(this.node, token);
    },
  },
  mixins: [bemMixin()],
  props: {
    node: {
      type: Node,
      required: true,
    }
  }
});

