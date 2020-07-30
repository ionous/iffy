// the node is a parameter of a parent run control
// it REPEATS
Vue.component('mk-repeater-ctrl', {
  template:
    `<span
      class="mk-repeater"
      :data-tag="node.type"
    ><template
      v-for="(kid, i) in node.kids"
      ><template v-if="commas"
        ><template v-if="mid(i)"
        >, </template
        ><template v-if="last(i)"
        > {{commas}} </template
      ></template
      ><mk-switch
        :node="kid"
        :key="kid.key"
      ></mk-switch
    ></template
  ></span>`,
  props: {
    node: Node,
  },
  // mixins: [bemMixin()],
  methods: {
    mid(i) {
      const { node }= this;
      return i && ((i>1) || !this.last(i));
    },
    last(i) {
      const { node }= this;
      return i === (node.kids.length - 1);
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
  computed: {
    commas() {
      const { node }= this;
      return node.kids.length > 1 && this.commaText(node.param.filters);
    },
  },
});
