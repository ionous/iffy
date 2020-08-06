
Vue.component('mk-opt-ctrl', {
  template:
  `<span
      :class="bemBlock()"
      :data-tag="node.type"
    ><mk-switch
      v-if="childNode"
      :node="childNode"
    ></mk-switch
    ><mk-pick-inline
      v-else
      :node="node"
      @picked="onPick"
    ></mk-pick-inline
  ></span>`,
  computed: {
    childNode() {
      return this.node.kid;
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
      const typeName= param.type || param; // an opt's param can map straight to their type.
      this.$root.redux.newSwap(node, token, typeName);
    },
  },
  mixins: [bemMixin()],
  props: {
    node: SwapNode,
    param: Object,
    token: String,
  }
});
