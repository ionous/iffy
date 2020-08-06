//
Vue.component('mk-repeater-ctrl', {
  template:
    `<span
      class="mk-repeater"
      :data-tag="node.type"
    ><template
      v-for="(kid, i) in nodes"
      ><template v-if="commas"
        ><template v-if="mid(i)"
        >, </template
        ><template v-if="last(i)"
        > {{commas}} </template
      ></template
      ><mk-switch
        :node="kid"
        :key="kid && kid.id"
        :param="param"
        :token="token"
      ></mk-switch
    ></template
  ></span>`,
  props: {
    node: Array, // node is really nodes for repeater control.
    param: Object,
    token: String,
  },
  computed: {
    nodes() {
      return this.node;
    },
    commas() {
      const { nodes, param }= this;
      return nodes.length > 1 && this.commaText(param.filters);
    },
  },

  methods: {
    mid(i) {
      return i && ((i>1) || !this.last(i));
    },
    last(i) {
      const { nodes }= this;
      return i === (nodes.length - 1);
    },
    commaText(filters) {
      let ret= "";
      if (filters) {
        if (filters.includes("comma-and")) {
          filters= "and";
        } else if (filters.includes("comma-or")) {
          filters= "or";
        }
      }
    },
  },
});
